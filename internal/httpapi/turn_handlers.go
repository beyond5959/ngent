package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/runtime"
	"github.com/beyond5959/ngent/internal/sse"
	"github.com/beyond5959/ngent/internal/storage"
)

func (s *Server) handleCreateTurnStream(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	if err := requireMethod(r, http.MethodPost); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	req, err := s.decodeTurnCreateRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "invalid request body", map[string]any{"reason": err.Error()})
		return
	}
	req.Prompt = agents.NormalizePrompt(req.Prompt)
	keepUploads := false
	defer func() {
		if keepUploads {
			return
		}
		removeStoredAttachments(req.Uploads)
	}()
	if !req.Stream {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "stream must be true", map[string]any{"field": "stream"})
		return
	}
	if len(req.Prompt.Content) == 0 {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "input or attachments are required", map[string]any{
			"fields": []string{"input", "attachments"},
		})
		return
	}

	injectedPrompt, err := s.buildInjectedPrompt(r.Context(), thread, req.Prompt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to build context window", map[string]any{
			"reason": err.Error(),
		})
		return
	}

	streamAgent, err := s.resolveTurnAgent(thread)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, codeUpstreamUnavailable, "failed to resolve agent provider", map[string]any{
			"agent":  thread.AgentID,
			"reason": err.Error(),
		})
		return
	}

	turnID := newTurnID()
	turnSessionID := threadSessionID(thread.AgentOptionsJSON)
	turnCtx, cancelTurn := context.WithCancel(context.WithoutCancel(r.Context()))
	persistCtx := context.WithoutCancel(r.Context())
	if err := s.turns.Activate(thread.ThreadID, turnSessionID, turnID, cancelTurn); err != nil {
		if errors.Is(err, runtime.ErrActiveTurnExists) {
			writeError(w, http.StatusConflict, codeConflict, "session already has an active turn", map[string]any{
				"threadId":  thread.ThreadID,
				"sessionId": turnSessionID,
			})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to activate turn", map[string]any{"reason": err.Error()})
		return
	}
	defer func() {
		cancelTurn()
		s.turns.Release(thread.ThreadID, turnSessionID, turnID)
		s.turnEvents.CloseTurn(turnID)
	}()
	if err := s.syncThreadConfigSelections(r.Context(), thread, streamAgent); err != nil {
		writeError(w, http.StatusServiceUnavailable, codeUpstreamUnavailable, "failed to sync thread config options", map[string]any{
			"threadId": thread.ThreadID,
			"reason":   err.Error(),
		})
		return
	}

	if _, err := s.store.CreateTurn(r.Context(), storage.CreateTurnParams{
		TurnID:      turnID,
		ThreadID:    thread.ThreadID,
		RequestText: req.Prompt.LegacyText(),
		Status:      "running",
		IsInternal:  false,
	}); err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to create turn", map[string]any{"reason": err.Error()})
		return
	}
	if err := s.persistTurnAttachments(persistCtx, turnID, req.Uploads); err != nil {
		s.finalizeTurnWithBestEffort(persistCtx, turnID, "failed", "error", "", err.Error())
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to persist turn attachments", map[string]any{
			"reason": err.Error(),
		})
		return
	}
	keepUploads = true

	streamWriter, err := sse.NewWriter(w)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "SSE is not supported by response writer", map[string]any{})
		return
	}

	aggregated := strings.Builder{}
	requestStreamClosed := false

	writeStreamEvent := func(eventType string, payload map[string]any) {
		if requestStreamClosed {
			return
		}
		if err := streamWriter.Event(eventType, payload); err != nil {
			requestStreamClosed = true
		}
	}

	emit := func(eventType string, payload map[string]any) error {
		_, streamPayload, err := s.appendTurnEvent(persistCtx, turnID, eventType, payload)
		if err != nil {
			return err
		}
		writeStreamEvent(eventType, streamPayload)
		return nil
	}
	appendOnlyEvent := func(eventType string, payload map[string]any) error {
		_, _, err := s.appendTurnEvent(persistCtx, turnID, eventType, payload)
		return err
	}

	if req.Prompt.HasResourceLinks() {
		if err := appendOnlyEvent(eventTypeUserPrompt, req.Prompt.EventPayload(turnID)); err != nil {
			s.finalizeTurnWithBestEffort(persistCtx, turnID, "failed", "error", "", err.Error())
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to persist user prompt", map[string]any{"reason": err.Error()})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	streamWriter.Flush()

	turnCtx = agents.WithPermissionHandler(turnCtx, func(permissionCtx context.Context, req agents.PermissionRequest) (agents.PermissionResponse, error) {
		permissionID := s.nextPermissionID(req.RequestID)
		pending := newPendingPermission(turnID, req.Options)
		s.registerPermission(permissionID, pending)
		defer s.unregisterPermission(permissionID, pending)

		payload := map[string]any{
			"turnId":       turnID,
			"permissionId": permissionID,
			"approval":     req.Approval,
			"command":      req.Command,
			"requestId":    req.RequestID,
		}
		if len(req.Options) > 0 {
			payload["options"] = req.Options
		}
		if err := emit("permission_required", payload); err != nil {
			pending.Resolve(permissionFailClosedResponse())
			return permissionFailClosedResponse(), err
		}

		response := s.waitPermissionResponse(permissionCtx, pending)
		if err := emit(eventTypePermissionResolved, pending.resolvedEventPayload(permissionID, response)); err != nil {
			s.logger.Warn("permission.resolved_event_failed",
				"turnId", turnID,
				"permissionId", permissionID,
				"reason", err.Error(),
			)
		}
		return response, nil
	})
	turnCtx = agents.WithPlanHandler(turnCtx, func(planCtx context.Context, entries []agents.PlanEntry) error {
		_ = planCtx
		payloadEntries := agents.ClonePlanEntries(entries)
		if payloadEntries == nil {
			payloadEntries = []agents.PlanEntry{}
		}
		return emit("plan_update", map[string]any{
			"turnId":  turnID,
			"entries": payloadEntries,
		})
	})
	turnCtx = agents.WithReasoningHandler(turnCtx, func(reasoningCtx context.Context, delta string) error {
		_ = reasoningCtx
		return emit(eventTypeReasoningDelta, map[string]any{
			"turnId": turnID,
			"delta":  delta,
		})
	})
	turnCtx = agents.WithSessionInfoHandler(turnCtx, func(sessionInfoCtx context.Context, update agents.SessionInfoUpdate) error {
		_ = sessionInfoCtx
		return emit(eventTypeSessionInfoUpdate, map[string]any{
			"turnId":    turnID,
			"sessionId": update.SessionID,
			"title":     update.Title,
		})
	})
	turnCtx = agents.WithSessionUsageHandler(turnCtx, func(sessionUsageCtx context.Context, update agents.SessionUsageUpdate) error {
		_ = sessionUsageCtx
		s.persistSessionUsageSnapshotBestEffort(persistCtx, thread, update)
		return emit(eventTypeSessionUsageUpdate, sessionUsageEventPayload(turnID, update))
	})
	turnCtx = agents.WithMessageContentHandler(turnCtx, func(messageCtx context.Context, event agents.ACPMessageContent) error {
		_ = messageCtx
		return emit(eventTypeMessageContent, event.EventPayload(turnID))
	})
	turnCtx = agents.WithToolCallHandler(turnCtx, func(toolCallCtx context.Context, event agents.ACPToolCall) error {
		_ = toolCallCtx
		eventType := strings.TrimSpace(event.Type)
		switch eventType {
		case eventTypeToolCall, eventTypeToolCallUpdate:
		default:
			return nil
		}
		return emit(eventType, event.EventPayload(turnID))
	})
	turnCtx = agents.WithSlashCommandsHandler(turnCtx, func(commandsCtx context.Context, commands []agents.SlashCommand) error {
		_ = commandsCtx
		if err := s.persistAgentSlashCommands(persistCtx, thread.AgentID, commands); err != nil {
			s.logger.Warn("thread.slash_commands_persist_failed",
				"threadId", thread.ThreadID,
				"agent", thread.AgentID,
				"reason", err.Error(),
			)
		}
		return nil
	})
	turnCtx = agents.WithConfigOptionsHandler(turnCtx, func(configOptionsCtx context.Context, options []agents.ConfigOption) error {
		_ = configOptionsCtx
		s.persistThreadConfigSnapshotBestEffort(persistCtx, &thread, options)
		return nil
	})
	turnCtx = agents.WithSessionBoundHandler(turnCtx, func(sessionCtx context.Context, sessionID string) error {
		_ = sessionCtx
		sessionID = strings.TrimSpace(sessionID)
		if sessionID == "" {
			return nil
		}
		if err := s.turns.BindTurnSession(turnID, sessionID); err != nil {
			if !errors.Is(err, runtime.ErrActiveTurnExists) {
				s.logger.Warn("thread.session_bind_runtime_failed",
					"threadId", thread.ThreadID,
					"agent", thread.AgentID,
					"turnId", turnID,
					"reason", err.Error(),
				)
			}
		} else {
			turnSessionID = sessionID
		}

		nextAgentOptionsJSON, changed, err := withThreadSessionID(thread.AgentOptionsJSON, sessionID)
		if err != nil {
			s.logger.Warn("thread.session_bind_encode_failed",
				"threadId", thread.ThreadID,
				"agent", thread.AgentID,
				"reason", err.Error(),
			)
		} else if changed {
			if err := s.store.UpdateThreadAgentOptions(persistCtx, thread.ThreadID, nextAgentOptionsJSON); err != nil {
				s.logger.Warn("thread.session_bind_persist_failed",
					"threadId", thread.ThreadID,
					"agent", thread.AgentID,
					"reason", err.Error(),
				)
			} else {
				if err := s.rebindManagedAgentScope(thread.ThreadID, thread.AgentOptionsJSON, nextAgentOptionsJSON); err != nil {
					s.logger.Warn("thread.session_bind_agent_scope_failed",
						"threadId", thread.ThreadID,
						"agent", thread.AgentID,
						"turnId", turnID,
						"reason", err.Error(),
					)
				}
				thread.AgentOptionsJSON = nextAgentOptionsJSON
				s.persistThreadSessionConfigSnapshotBestEffort(persistCtx, thread)
			}
		}

		if err := emit("session_bound", map[string]any{
			"threadId":  thread.ThreadID,
			"sessionId": sessionID,
		}); err != nil {
			s.logger.Warn("thread.session_bind_emit_failed",
				"threadId", thread.ThreadID,
				"agent", thread.AgentID,
				"reason", err.Error(),
			)
		}
		return nil
	})

	if err := emit("turn_started", map[string]any{"turnId": turnID}); err != nil {
		s.finalizeTurnWithBestEffort(persistCtx, turnID, "failed", "error", "", err.Error())
		return
	}

	stopReason, streamErr := agents.StreamPrompt(turnCtx, streamAgent, injectedPrompt, func(delta string) error {
		aggregated.WriteString(delta)
		return emit("message_delta", map[string]any{"turnId": turnID, "delta": delta})
	})

	finalStatus := "completed"
	finalReason := string(agents.StopReasonEndTurn)
	errorMessage := ""

	if streamErr != nil {
		finalStatus = "failed"
		finalReason = "error"
		errorMessage = streamErr.Error()
		_ = emit("error", map[string]any{
			"turnId":  turnID,
			"code":    classifyStreamErrorCode(streamErr),
			"message": streamErr.Error(),
		})
	} else if stopReason == agents.StopReasonCancelled {
		finalStatus = "cancelled"
		finalReason = string(agents.StopReasonCancelled)
	}

	if err := emit("turn_completed", map[string]any{"turnId": turnID, "stopReason": finalReason}); err != nil && errorMessage == "" {
		errorMessage = err.Error()
		if finalStatus == "completed" {
			finalStatus = "failed"
			finalReason = "error"
		}
	}

	s.finalizeTurnWithBestEffort(persistCtx, turnID, finalStatus, finalReason, aggregated.String(), errorMessage)
}

func (s *Server) persistTurnAttachments(ctx context.Context, turnID string, uploads []storedTurnAttachment) error {
	if len(uploads) == 0 {
		return nil
	}

	params := make([]storage.CreateTurnAttachmentParams, 0, len(uploads))
	for _, upload := range uploads {
		attachmentID := strings.TrimSpace(upload.PromptContent.AttachmentID)
		if attachmentID == "" {
			return errors.New("attachmentID is required")
		}
		params = append(params, storage.CreateTurnAttachmentParams{
			AttachmentID: attachmentID,
			TurnID:       turnID,
			Name:         strings.TrimSpace(upload.PromptContent.Name),
			MimeType:     strings.TrimSpace(upload.PromptContent.MimeType),
			Size:         upload.PromptContent.Size,
			FilePath:     strings.TrimSpace(upload.FilePath),
		})
	}
	return s.store.CreateTurnAttachments(ctx, params)
}

func (s *Server) handleCompactThread(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	if err := requireMethod(r, http.MethodPost); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	var req struct {
		MaxSummaryChars int `json:"maxSummaryChars"`
	}
	if r.Body != nil {
		if err := decodeJSONBody(r, &req); err != nil && !errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "invalid JSON body", map[string]any{"reason": err.Error()})
			return
		}
	}

	summaryLimit := req.MaxSummaryChars
	if summaryLimit <= 0 {
		summaryLimit = s.compactMaxChars
	}

	streamAgent, err := s.resolveTurnAgent(thread)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, codeUpstreamUnavailable, "failed to resolve agent provider", map[string]any{
			"agent":  thread.AgentID,
			"reason": err.Error(),
		})
		return
	}

	compactPrompt, err := s.buildCompactPrompt(r.Context(), thread, summaryLimit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to build compact prompt", map[string]any{
			"reason": err.Error(),
		})
		return
	}

	turnID := newTurnID()
	turnCtx, cancelTurn := context.WithCancel(r.Context())
	persistCtx := context.WithoutCancel(r.Context())
	if err := s.turns.ActivateThreadExclusive(thread.ThreadID, turnID, cancelTurn); err != nil {
		if errors.Is(err, runtime.ErrActiveTurnExists) {
			writeError(w, http.StatusConflict, codeConflict, "thread already has an active turn", map[string]any{"threadId": thread.ThreadID})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to activate compact turn", map[string]any{"reason": err.Error()})
		return
	}
	defer func() {
		cancelTurn()
		s.turns.ReleaseThreadExclusive(thread.ThreadID, turnID)
	}()

	if _, err := s.store.CreateTurn(r.Context(), storage.CreateTurnParams{
		TurnID:      turnID,
		ThreadID:    thread.ThreadID,
		RequestText: compactPrompt,
		Status:      "running",
		IsInternal:  true,
	}); err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to create compact turn", map[string]any{"reason": err.Error()})
		return
	}

	appendOnlyEvent := func(eventType string, payload map[string]any) error {
		dataJSON, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return marshalErr
		}
		_, appendErr := s.store.AppendEvent(persistCtx, turnID, eventType, string(dataJSON))
		return appendErr
	}

	if err := appendOnlyEvent("turn_started", map[string]any{"turnId": turnID}); err != nil {
		s.finalizeTurnWithBestEffort(persistCtx, turnID, "failed", "error", "", err.Error())
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to persist compact start event", map[string]any{"reason": err.Error()})
		return
	}

	aggregated := strings.Builder{}
	turnCtx = agents.WithPlanHandler(turnCtx, func(planCtx context.Context, entries []agents.PlanEntry) error {
		_ = planCtx
		payloadEntries := agents.ClonePlanEntries(entries)
		if payloadEntries == nil {
			payloadEntries = []agents.PlanEntry{}
		}
		return appendOnlyEvent("plan_update", map[string]any{
			"turnId":  turnID,
			"entries": payloadEntries,
		})
	})
	turnCtx = agents.WithReasoningHandler(turnCtx, func(reasoningCtx context.Context, delta string) error {
		_ = reasoningCtx
		return appendOnlyEvent(eventTypeReasoningDelta, map[string]any{
			"turnId": turnID,
			"delta":  delta,
		})
	})
	turnCtx = agents.WithSessionUsageHandler(turnCtx, func(sessionUsageCtx context.Context, update agents.SessionUsageUpdate) error {
		_ = sessionUsageCtx
		s.persistSessionUsageSnapshotBestEffort(persistCtx, thread, update)
		return appendOnlyEvent(eventTypeSessionUsageUpdate, sessionUsageEventPayload(turnID, update))
	})
	turnCtx = agents.WithMessageContentHandler(turnCtx, func(messageCtx context.Context, event agents.ACPMessageContent) error {
		_ = messageCtx
		return appendOnlyEvent(eventTypeMessageContent, event.EventPayload(turnID))
	})
	turnCtx = agents.WithToolCallHandler(turnCtx, func(toolCallCtx context.Context, event agents.ACPToolCall) error {
		_ = toolCallCtx
		eventType := strings.TrimSpace(event.Type)
		switch eventType {
		case eventTypeToolCall, eventTypeToolCallUpdate:
		default:
			return nil
		}
		return appendOnlyEvent(eventType, event.EventPayload(turnID))
	})
	stopReason, streamErr := streamAgent.Stream(turnCtx, compactPrompt, func(delta string) error {
		aggregated.WriteString(delta)
		return appendOnlyEvent("message_delta", map[string]any{
			"turnId": turnID,
			"delta":  delta,
		})
	})

	finalStatus := "completed"
	finalReason := string(agents.StopReasonEndTurn)
	errorMessage := ""

	if streamErr != nil {
		finalStatus = "failed"
		finalReason = "error"
		errorMessage = streamErr.Error()
		_ = appendOnlyEvent("error", map[string]any{
			"turnId":  turnID,
			"code":    classifyStreamErrorCode(streamErr),
			"message": streamErr.Error(),
		})
	} else if stopReason == agents.StopReasonCancelled {
		finalStatus = "cancelled"
		finalReason = string(agents.StopReasonCancelled)
	}

	if err := appendOnlyEvent("turn_completed", map[string]any{"turnId": turnID, "stopReason": finalReason}); err != nil && errorMessage == "" {
		errorMessage = err.Error()
		if finalStatus == "completed" {
			finalStatus = "failed"
			finalReason = "error"
		}
	}

	newSummary := clampToChars(strings.TrimSpace(aggregated.String()), summaryLimit)
	if finalStatus == "completed" && finalReason == string(agents.StopReasonEndTurn) {
		if err := s.store.UpdateThreadSummary(persistCtx, thread.ThreadID, newSummary); err != nil {
			finalStatus = "failed"
			finalReason = "error"
			errorMessage = err.Error()
		}
	}

	s.finalizeTurnWithBestEffort(persistCtx, turnID, finalStatus, finalReason, aggregated.String(), errorMessage)

	if finalStatus != "completed" {
		statusCode := http.StatusInternalServerError
		errorCode := codeInternal
		if streamErr != nil {
			errorCode = classifyStreamErrorCode(streamErr)
			switch errorCode {
			case codeTimeout:
				statusCode = http.StatusGatewayTimeout
			case codeUpstreamUnavailable:
				statusCode = http.StatusServiceUnavailable
			}
		}
		writeError(w, statusCode, errorCode, "compact failed", map[string]any{
			"turnId": turnID,
			"reason": errorMessage,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"threadId":     thread.ThreadID,
		"turnId":       turnID,
		"status":       finalStatus,
		"stopReason":   finalReason,
		"summary":      newSummary,
		"summaryChars": runeLen(newSummary),
	})
}

func (s *Server) handleCancelTurn(w http.ResponseWriter, r *http.Request, clientID, turnID string) {
	if err := requireMethod(r, http.MethodPost); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	turn, err := s.store.GetTurn(r.Context(), turnID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, codeNotFound, "turn not found", map[string]any{})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to load turn", map[string]any{"reason": err.Error()})
		return
	}

	thread, ok := s.getAccessibleThread(r.Context(), turn.ThreadID)
	if !ok {
		writeError(w, http.StatusNotFound, codeNotFound, "turn not found", map[string]any{})
		return
	}

	if err := s.turns.Cancel(turnID); err != nil {
		if errors.Is(err, runtime.ErrTurnNotActive) {
			writeError(w, http.StatusConflict, codeConflict, "turn is not active", map[string]any{"turnId": turnID})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to cancel turn", map[string]any{"reason": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"turnId":   turnID,
		"threadId": thread.ThreadID,
		"status":   "cancelling",
	})
}

func (s *Server) finalizeTurnWithBestEffort(ctx context.Context, turnID, status, stopReason, responseText, errorMessage string) {
	_ = s.store.FinalizeTurn(ctx, storage.FinalizeTurnParams{
		TurnID:       turnID,
		ResponseText: responseText,
		Status:       status,
		StopReason:   stopReason,
		ErrorMessage: errorMessage,
	})
}

func (s *Server) buildInjectedPrompt(ctx context.Context, thread storage.Thread, prompt agents.Prompt) (agents.Prompt, error) {
	prompt = agents.NormalizePrompt(prompt)
	if threadSessionID(thread.AgentOptionsJSON) != "" || threadFreshSessionRequested(thread.AgentOptionsJSON) {
		return prompt, nil
	}

	recentTurns, err := s.loadRecentVisibleTurns(ctx, thread.ThreadID)
	if err != nil {
		return agents.Prompt{}, err
	}

	currentInput := prompt.Text()
	if strings.TrimSpace(thread.Summary) == "" && len(recentTurns) == 0 && currentInput == "" {
		return prompt, nil
	}

	content := make([]agents.PromptContent, 0, len(prompt.Content))
	injectedText := composeContextPrompt(
		thread.Summary,
		recentTurns,
		currentInput,
		s.contextMaxChars,
	)
	if strings.TrimSpace(injectedText) != "" {
		content = append(content, agents.PromptContent{
			Type: agents.PromptContentTypeText,
			Text: injectedText,
		})
	}
	for _, item := range prompt.Content {
		if item.Type == agents.PromptContentTypeResourceLink {
			content = append(content, item)
		}
	}
	return agents.NormalizePrompt(agents.Prompt{Content: content}), nil
}

func (s *Server) buildCompactPrompt(ctx context.Context, thread storage.Thread, maxSummaryChars int) (string, error) {
	recentTurns, err := s.loadRecentVisibleTurns(ctx, thread.ThreadID)
	if err != nil {
		return "", err
	}

	instruction := fmt.Sprintf(
		"Please generate an updated rolling summary of the conversation. "+
			"Output plain text only, keep key decisions/constraints, and limit to %d characters.",
		maxSummaryChars,
	)
	return composeContextPrompt(
		thread.Summary,
		recentTurns,
		instruction,
		s.contextMaxChars,
	), nil
}

func (s *Server) loadRecentVisibleTurns(ctx context.Context, threadID string) ([]storage.Turn, error) {
	turns, err := s.store.ListTurnsByThread(ctx, threadID)
	if err != nil {
		return nil, err
	}

	filtered := make([]storage.Turn, 0, len(turns))
	for _, turn := range turns {
		if turn.IsInternal {
			continue
		}
		filtered = append(filtered, turn)
	}

	if len(filtered) > s.contextRecentTurns {
		filtered = filtered[len(filtered)-s.contextRecentTurns:]
	}
	return filtered, nil
}

func composeContextPrompt(summary string, recentTurns []storage.Turn, currentInput string, maxChars int) string {
	summary = strings.TrimSpace(summary)
	currentInput = strings.TrimSpace(currentInput)

	recentCopy := make([]storage.Turn, len(recentTurns))
	copy(recentCopy, recentTurns)

	if summary == "" && len(recentCopy) == 0 {
		if maxChars <= 0 || runeLen(currentInput) <= maxChars {
			return currentInput
		}
		return clampToChars(currentInput, maxChars)
	}

	for i := 0; i < 256; i++ {
		prompt := renderContextPrompt(summary, recentCopy, currentInput)
		if maxChars <= 0 || runeLen(prompt) <= maxChars {
			return prompt
		}

		if len(recentCopy) > 0 {
			recentCopy = recentCopy[1:]
			continue
		}

		if runeLen(summary) > 0 {
			summary = clampToChars(summary, runeLen(summary)-maxInt(1, runeLen(summary)/4))
			continue
		}

		if runeLen(currentInput) > 0 {
			currentInput = truncateFromEnd(currentInput, runeLen(currentInput)-maxInt(1, runeLen(currentInput)/4))
			continue
		}

		return clampToChars(prompt, maxChars)
	}

	return clampToChars(renderContextPrompt(summary, recentCopy, currentInput), maxChars)
}

func renderContextPrompt(summary string, recentTurns []storage.Turn, currentInput string) string {
	var builder strings.Builder
	builder.WriteString("[Conversation Summary]\n")
	if summary == "" {
		builder.WriteString("(empty)")
	} else {
		builder.WriteString(summary)
	}

	builder.WriteString("\n\n[Recent Turns]\n")
	if len(recentTurns) == 0 {
		builder.WriteString("(none)")
	} else {
		for _, turn := range recentTurns {
			builder.WriteString("User: ")
			builder.WriteString(strings.TrimSpace(turn.RequestText))
			builder.WriteString("\nAssistant: ")
			builder.WriteString(strings.TrimSpace(turn.ResponseText))
			builder.WriteString("\n")
		}
		builder.WriteString("----")
	}

	builder.WriteString("\n\n[Current User Input]\n")
	builder.WriteString(currentInput)
	return builder.String()
}

func clampToChars(text string, maxChars int) string {
	if maxChars <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= maxChars {
		return text
	}
	return string(runes[:maxChars])
}

func truncateFromEnd(text string, maxChars int) string {
	if maxChars <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= maxChars {
		return text
	}
	return string(runes[len(runes)-maxChars:])
}

func runeLen(text string) int {
	return len([]rune(text))
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func classifyStreamErrorCode(err error) string {
	if err == nil {
		return codeInternal
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return codeTimeout
	}
	if errors.Is(err, context.Canceled) {
		return codeTimeout
	}
	return codeUpstreamUnavailable
}
