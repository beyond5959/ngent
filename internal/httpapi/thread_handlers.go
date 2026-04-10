package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/runtime"
	"github.com/beyond5959/ngent/internal/storage"
)

type threadResponse struct {
	ThreadID         string          `json:"threadId"`
	Agent            string          `json:"agent"`
	CWD              string          `json:"cwd"`
	Title            string          `json:"title"`
	AgentOptions     json.RawMessage `json:"agentOptions"`
	Summary          string          `json:"summary"`
	HasActiveSession bool            `json:"hasActiveSession"`
	CreatedAt        string          `json:"createdAt"`
	UpdatedAt        string          `json:"updatedAt"`
}

type turnHistoryResponse struct {
	TurnID       string                 `json:"turnId"`
	RequestText  string                 `json:"requestText"`
	ResponseText string                 `json:"responseText"`
	IsInternal   bool                   `json:"isInternal,omitempty"`
	Status       string                 `json:"status"`
	StopReason   string                 `json:"stopReason"`
	ErrorMessage string                 `json:"errorMessage"`
	CreatedAt    string                 `json:"createdAt"`
	CompletedAt  *string                `json:"completedAt,omitempty"`
	Events       []eventHistoryResponse `json:"events,omitempty"`
}

type eventHistoryResponse struct {
	EventID   int64           `json:"eventId"`
	Seq       int             `json:"seq"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	CreatedAt string          `json:"createdAt"`
}

type historySessionTranscriptResponse struct {
	Supported bool                              `json:"supported"`
	Messages  []agents.SessionTranscriptMessage `json:"messages"`
}

type threadHistoryTurn struct {
	turn   storage.Turn
	events []storage.Event
}

type threadHistoryTurnAssignment struct {
	turn      threadHistoryTurn
	sessionID string
}

func (s *Server) handleCreateThread(w http.ResponseWriter, r *http.Request, clientID string) {
	var req struct {
		Agent        string          `json:"agent"`
		CWD          string          `json:"cwd"`
		Title        string          `json:"title"`
		AgentOptions json.RawMessage `json:"agentOptions"`
	}

	if err := requireMethod(r, http.MethodPost); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	if err := decodeJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "invalid JSON body", map[string]any{"reason": err.Error()})
		return
	}

	req.Agent = strings.TrimSpace(req.Agent)
	if _, ok := s.allowedAgent[req.Agent]; !ok {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "agent is not in allowlist", map[string]any{
			"field":         "agent",
			"allowedAgents": sortedAgentIDs(s.allowedAgent),
		})
		return
	}

	cwd := strings.TrimSpace(req.CWD)
	if cwd == "" {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "cwd is required", map[string]any{"field": "cwd"})
		return
	}

	expandedCWD, err := expandPath(cwd)
	if err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "failed to expand path", map[string]any{"field": "cwd", "reason": err.Error()})
		return
	}
	cwd = expandedCWD

	if !filepath.IsAbs(cwd) {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "cwd must be an absolute path", map[string]any{"field": "cwd"})
		return
	}
	cwd = filepath.Clean(cwd)
	if !isPathAllowed(cwd, s.allowedRoots) {
		writeError(w, http.StatusForbidden, codeForbidden, "cwd is outside allowed roots", map[string]any{
			"field":         "cwd",
			"cwd":           cwd,
			"allowed_roots": s.allowedRoots,
		})
		return
	}

	if _, err := os.Stat(cwd); err != nil {
		if os.IsNotExist(err) {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "cwd does not exist", map[string]any{
				"field": "cwd",
				"cwd":   cwd,
			})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to check cwd", map[string]any{
			"field":  "cwd",
			"reason": err.Error(),
		})
		return
	}

	agentOptionsJSON, err := normalizeAgentOptions(req.AgentOptions)
	if err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "agentOptions must be a JSON object", map[string]any{"field": "agentOptions"})
		return
	}

	threadID := newThreadID()
	_, err = s.store.CreateThread(r.Context(), storage.CreateThreadParams{
		ThreadID:         threadID,
		AgentID:          req.Agent,
		CWD:              cwd,
		Title:            req.Title,
		AgentOptionsJSON: agentOptionsJSON,
		Summary:          "",
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to create thread", map[string]any{"reason": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"threadId": threadID})
}

func (s *Server) handleListThreads(w http.ResponseWriter, r *http.Request, clientID string) {
	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	threads, err := s.store.ListThreads(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to list threads", map[string]any{"reason": err.Error()})
		return
	}

	items := make([]threadResponse, 0, len(threads))
	for _, thread := range threads {
		item, convErr := s.threadResponseForThread(thread)
		if convErr != nil {
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to encode thread", map[string]any{"reason": convErr.Error()})
			return
		}
		items = append(items, item)
	}

	writeJSON(w, http.StatusOK, map[string]any{"threads": items})
}

func (s *Server) handleGetThread(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	resp, convErr := s.threadResponseForThread(thread)
	if convErr != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to encode thread", map[string]any{"reason": convErr.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"thread": resp})
}

func (s *Server) handleUpdateThread(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	if err := requireMethod(r, http.MethodPatch); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	var req struct {
		Title        *string          `json:"title"`
		AgentOptions *json.RawMessage `json:"agentOptions"`
	}
	if err := decodeJSONBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "invalid JSON body", map[string]any{"reason": err.Error()})
		return
	}

	agentOptionsJSON := ""
	sessionOnlyUpdate := false
	currentSessionID := threadSessionID(thread.AgentOptionsJSON)
	currentFreshSession := threadFreshSessionRequested(thread.AgentOptionsJSON)
	if req.AgentOptions != nil {
		var err error
		agentOptionsJSON, err = normalizeAgentOptions(*req.AgentOptions)
		if err != nil {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "agentOptions must be a JSON object", map[string]any{"field": "agentOptions"})
			return
		}

		nextSessionID := threadSessionID(agentOptionsJSON)
		if nextSessionID != "" && nextSessionID != currentSessionID {
			agentOptionsJSON, err = withoutThreadConfigState(agentOptionsJSON)
			if err != nil {
				writeError(w, http.StatusBadRequest, codeInvalidArgument, "agentOptions must be a JSON object", map[string]any{"field": "agentOptions"})
				return
			}
		}
		sessionOnlyUpdate, err = isSessionOnlyAgentOptionsUpdate(thread.AgentOptionsJSON, agentOptionsJSON)
		if err != nil {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "agentOptions must be a JSON object", map[string]any{"field": "agentOptions"})
			return
		}

		nextSessionID = threadSessionID(agentOptionsJSON)
		shouldRequestFreshSession := currentFreshSession
		switch {
		case nextSessionID != "":
			shouldRequestFreshSession = false
		case currentSessionID != "":
			shouldRequestFreshSession = true
		}

		agentOptionsJSON, _, err = withThreadFreshSessionRequested(agentOptionsJSON, shouldRequestFreshSession)
		if err != nil {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "agentOptions must be a JSON object", map[string]any{"field": "agentOptions"})
			return
		}
	}

	allowSessionSelectionWhileActive := req.Title == nil && req.AgentOptions != nil && sessionOnlyUpdate
	if s.turns.IsThreadActive(thread.ThreadID) && !allowSessionSelectionWhileActive {
		writeError(w, http.StatusConflict, codeConflict, "thread has an active turn", map[string]any{"threadId": thread.ThreadID})
		return
	}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if err := s.store.UpdateThreadTitle(r.Context(), thread.ThreadID, title); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				writeThreadNotFound(w)
				return
			}
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to update thread", map[string]any{"reason": err.Error()})
			return
		}
	}

	if req.AgentOptions != nil {
		if err := s.store.UpdateThreadAgentOptions(r.Context(), thread.ThreadID, agentOptionsJSON); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				writeThreadNotFound(w)
				return
			}
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to update thread", map[string]any{"reason": err.Error()})
			return
		}

		nextSessionID := threadSessionID(agentOptionsJSON)
		if sessionOnlyUpdate && currentSessionID != nextSessionID && nextSessionID == "" {
			s.closeThreadAgentScope(thread.ThreadID, agentOptionsJSON, "thread_session_reset")
			if plainAgentOptionsJSON, changed, err := withThreadFreshSessionRequested(agentOptionsJSON, false); err == nil && changed {
				s.closeThreadAgentScope(thread.ThreadID, plainAgentOptionsJSON, "thread_session_reset")
			}
		}
		if !sessionOnlyUpdate {
			s.closeThreadAgents(thread.ThreadID, "thread_updated")
		}
	}

	updatedThread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), thread.ThreadID)
	if !ok {
		return
	}

	resp, convErr := s.threadResponseForThread(updatedThread)
	if convErr != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to encode thread", map[string]any{"reason": convErr.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"thread": resp})
}

func (s *Server) handleDeleteThread(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	if err := requireMethod(r, http.MethodDelete); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	deleteGuardTurnID := "delete-" + newTurnID()
	if err := s.turns.ActivateThreadExclusive(thread.ThreadID, deleteGuardTurnID, nil); err != nil {
		if errors.Is(err, runtime.ErrActiveTurnExists) {
			writeError(w, http.StatusConflict, codeConflict, "thread has an active turn", map[string]any{"threadId": thread.ThreadID})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to lock thread for delete", map[string]any{"reason": err.Error()})
		return
	}
	defer s.turns.ReleaseThreadExclusive(thread.ThreadID, deleteGuardTurnID)

	if err := s.store.DeleteThread(r.Context(), thread.ThreadID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeThreadNotFound(w)
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to delete thread", map[string]any{"reason": err.Error()})
		return
	}

	s.closeThreadAgents(thread.ThreadID, "thread_deleted")

	writeJSON(w, http.StatusOK, map[string]any{
		"threadId": thread.ThreadID,
		"status":   "deleted",
	})
}

func (s *Server) handleThreadHistory(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	includeEvents := parseBoolQuery(r, "includeEvents")
	includeInternal := parseBoolQuery(r, "includeInternal")
	sessionID := strings.TrimSpace(r.URL.Query().Get("sessionId"))

	turns, err := s.store.ListTurnsByThread(r.Context(), threadID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to list history", map[string]any{"reason": err.Error()})
		return
	}

	loadEvents := includeEvents || sessionID != ""
	eventsByTurnID := map[string][]storage.Event(nil)
	if loadEvents {
		turnIDs := make([]string, 0, len(turns))
		for _, turn := range turns {
			turnIDs = append(turnIDs, turn.TurnID)
		}
		eventsByTurnID, err = s.store.ListEventsByTurns(r.Context(), turnIDs)
		if err != nil {
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to list events", map[string]any{"reason": err.Error()})
			return
		}
	}

	historyTurns := make([]threadHistoryTurn, 0, len(turns))
	for _, turn := range turns {
		historyTurn := threadHistoryTurn{turn: turn}
		if loadEvents {
			historyTurn.events = eventsByTurnID[turn.TurnID]
		}
		if !includeInternal && turn.IsInternal {
			continue
		}
		historyTurns = append(historyTurns, historyTurn)
	}

	if sessionID != "" {
		historyTurns = filterThreadHistoryBySession(historyTurns, sessionID)
	}

	respTurns := make([]turnHistoryResponse, 0, len(historyTurns))
	for _, item := range historyTurns {
		turn := item.turn

		respTurn := turnHistoryResponse{
			TurnID:       turn.TurnID,
			RequestText:  turn.RequestText,
			ResponseText: turn.ResponseText,
			IsInternal:   turn.IsInternal,
			Status:       turn.Status,
			StopReason:   turn.StopReason,
			ErrorMessage: turn.ErrorMessage,
			CreatedAt:    turn.CreatedAt.UTC().Format(time.RFC3339Nano),
		}
		if turn.CompletedAt != nil {
			completed := turn.CompletedAt.UTC().Format(time.RFC3339Nano)
			respTurn.CompletedAt = &completed
		}

		if includeEvents {
			events := compactThreadHistoryEvents(item.events)
			respEvents := make([]eventHistoryResponse, 0, len(events))
			for _, event := range events {
				raw := json.RawMessage(event.DataJSON)
				if len(strings.TrimSpace(event.DataJSON)) == 0 || !json.Valid(raw) {
					raw = json.RawMessage(`{}`)
				}
				respEvents = append(respEvents, eventHistoryResponse{
					EventID:   event.EventID,
					Seq:       event.Seq,
					Type:      event.Type,
					Data:      raw,
					CreatedAt: event.CreatedAt.UTC().Format(time.RFC3339Nano),
				})
			}
			respTurn.Events = respEvents
		}

		respTurns = append(respTurns, respTurn)
	}

	resp := map[string]any{"turns": respTurns}
	if sessionID != "" {
		transcript, transcriptErr := s.loadThreadSessionTranscript(r.Context(), thread, sessionID, len(historyTurns) == 0)
		if transcriptErr != nil {
			s.logger.Warn("thread_history.session_transcript_load_failed",
				"threadId", thread.ThreadID,
				"agent", thread.AgentID,
				"sessionId", sessionID,
				"reason", transcriptErr.Error(),
			)
		} else if transcript != nil {
			resp["sessionTranscript"] = historySessionTranscriptResponse{
				Supported: transcript.Supported,
				Messages:  transcript.Messages,
			}
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

func filterThreadHistoryBySession(turns []threadHistoryTurn, sessionID string) []threadHistoryTurn {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return turns
	}

	assignments := make([]threadHistoryTurnAssignment, 0, len(turns))
	annotatedSessions := make(map[string]struct{})
	for _, turn := range turns {
		assignedSessionID := extractThreadHistorySessionID(turn.events)
		if assignedSessionID != "" {
			annotatedSessions[assignedSessionID] = struct{}{}
		}
		assignments = append(assignments, threadHistoryTurnAssignment{
			turn:      turn,
			sessionID: assignedSessionID,
		})
	}

	if len(annotatedSessions) == 0 {
		filtered := make([]threadHistoryTurn, 0, len(assignments))
		for _, item := range assignments {
			if isEphemeralCancelledHistoryTurn(item.turn.turn) {
				continue
			}
			filtered = append(filtered, item.turn)
		}
		return filtered
	}

	hasMatchedAnnotatedTurns := false
	for _, item := range assignments {
		if item.sessionID == sessionID {
			hasMatchedAnnotatedTurns = true
			break
		}
	}
	if !hasMatchedAnnotatedTurns {
		return nil
	}

	includeUnannotatedLegacyTurns := len(annotatedSessions) == 1
	if includeUnannotatedLegacyTurns {
		_, includeUnannotatedLegacyTurns = annotatedSessions[sessionID]
	}

	filtered := make([]threadHistoryTurn, 0, len(assignments))
	for _, item := range assignments {
		if item.sessionID == sessionID || (includeUnannotatedLegacyTurns && item.sessionID == "") {
			filtered = append(filtered, item.turn)
		}
	}
	return filtered
}

func extractThreadHistorySessionID(events []storage.Event) string {
	sessionID := ""
	for _, event := range events {
		if event.Type != "session_bound" {
			continue
		}
		var payload struct {
			SessionID string `json:"sessionId"`
		}
		if err := json.Unmarshal([]byte(event.DataJSON), &payload); err != nil {
			continue
		}
		nextSessionID := strings.TrimSpace(payload.SessionID)
		if nextSessionID != "" {
			sessionID = nextSessionID
		}
	}
	return sessionID
}

func isEphemeralCancelledHistoryTurn(turn storage.Turn) bool {
	return turn.Status == "cancelled" && strings.TrimSpace(turn.ResponseText) == ""
}

func compactThreadHistoryEvents(events []storage.Event) []storage.Event {
	if len(events) < 2 {
		return events
	}

	compacted := make([]storage.Event, 0, len(events))
	for _, event := range events {
		if lastIdx := len(compacted) - 1; lastIdx >= 0 {
			last := compacted[lastIdx]
			if shouldCompactThreadHistoryDeltaEvent(last, event) {
				mergedDataJSON, merged, err := compactThreadHistoryDeltaPayload(
					last.TurnID,
					last.DataJSON,
					event.DataJSON,
				)
				if err != nil {
					return events
				}
				if merged {
					compacted[lastIdx].DataJSON = mergedDataJSON
					compacted[lastIdx].CreatedAt = event.CreatedAt
					continue
				}
			}
		}

		compacted = append(compacted, event)
	}
	return compacted
}

func shouldCompactThreadHistoryDeltaEvent(last, next storage.Event) bool {
	if strings.TrimSpace(last.TurnID) != strings.TrimSpace(next.TurnID) {
		return false
	}
	if strings.TrimSpace(last.Type) != strings.TrimSpace(next.Type) {
		return false
	}

	switch strings.TrimSpace(next.Type) {
	case "message_delta", "reasoning_delta", "thought_delta":
		return true
	default:
		return false
	}
}

func compactThreadHistoryDeltaPayload(turnID, currentDataJSON, nextDataJSON string) (string, bool, error) {
	currentPayload := map[string]any{}
	if err := json.Unmarshal([]byte(currentDataJSON), &currentPayload); err != nil {
		return "", false, nil
	}
	nextPayload := map[string]any{}
	if err := json.Unmarshal([]byte(nextDataJSON), &nextPayload); err != nil {
		return "", false, nil
	}

	currentDelta, currentOK := currentPayload["delta"].(string)
	nextDelta, nextOK := nextPayload["delta"].(string)
	if !currentOK || !nextOK {
		return "", false, nil
	}
	if !threadHistoryDeltaPayloadMatchesTurn(turnID, currentPayload) {
		return "", false, nil
	}
	if !threadHistoryDeltaPayloadMatchesTurn(turnID, nextPayload) {
		return "", false, nil
	}

	currentPayload["delta"] = currentDelta + nextDelta
	mergedJSON, err := json.Marshal(currentPayload)
	if err != nil {
		return "", false, err
	}
	return string(mergedJSON), true, nil
}

func threadHistoryDeltaPayloadMatchesTurn(turnID string, payload map[string]any) bool {
	value, ok := payload["turnId"]
	if !ok {
		return true
	}
	valueText, ok := value.(string)
	if !ok {
		return false
	}
	return strings.TrimSpace(valueText) == strings.TrimSpace(turnID)
}

func (s *Server) getAccessibleThread(ctx context.Context, threadID string) (storage.Thread, bool) {
	thread, err := s.store.GetThread(ctx, threadID)
	if err != nil {
		return storage.Thread{}, false
	}
	return thread, true
}

func (s *Server) loadAccessibleThreadOrWriteNotFound(w http.ResponseWriter, ctx context.Context, threadID string) (storage.Thread, bool) {
	thread, ok := s.getAccessibleThread(ctx, threadID)
	if !ok {
		writeThreadNotFound(w)
		return storage.Thread{}, false
	}
	return thread, true
}

func writeThreadNotFound(w http.ResponseWriter) {
	writeError(w, http.StatusNotFound, codeNotFound, "thread not found", map[string]any{})
}

func (s *Server) threadResponseForThread(thread storage.Thread) (threadResponse, error) {
	return toThreadResponse(thread, s.turns.HasActiveSession(thread.ThreadID))
}

func toThreadResponse(thread storage.Thread, hasActiveSession bool) (threadResponse, error) {
	raw, err := sanitizeThreadAgentOptionsForResponse(thread.AgentOptionsJSON)
	if err != nil {
		return threadResponse{}, fmt.Errorf("invalid agent_options_json for thread %s", thread.ThreadID)
	}

	return threadResponse{
		ThreadID:         thread.ThreadID,
		Agent:            thread.AgentID,
		CWD:              thread.CWD,
		Title:            thread.Title,
		AgentOptions:     raw,
		Summary:          thread.Summary,
		HasActiveSession: hasActiveSession,
		CreatedAt:        thread.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:        thread.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}, nil
}
