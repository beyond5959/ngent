package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/storage"
)

type sessionInfoResponse struct {
	SessionID string         `json:"sessionId"`
	CWD       string         `json:"cwd,omitempty"`
	Title     string         `json:"title,omitempty"`
	UpdatedAt string         `json:"updatedAt,omitempty"`
	Meta      map[string]any `json:"_meta,omitempty"`
	IsActive  bool           `json:"isActive,omitempty"`
}

func mergeSessionInfo(current, incoming agents.SessionInfo) agents.SessionInfo {
	cloned := agents.CloneSessionInfo(current)
	next := agents.CloneSessionInfo(incoming)
	if next.SessionID != "" {
		cloned.SessionID = next.SessionID
	}
	if next.CWD != "" {
		cloned.CWD = next.CWD
	}
	if next.Title != "" {
		cloned.Title = next.Title
	}
	if next.UpdatedAt != "" {
		cloned.UpdatedAt = next.UpdatedAt
	}
	if len(next.Meta) > 0 {
		if cloned.Meta == nil {
			cloned.Meta = map[string]any{}
		}
		for key, value := range next.Meta {
			cloned.Meta[key] = value
		}
	}
	return cloned
}

func includeCurrentThreadSession(thread storage.Thread, result agents.SessionListResult) agents.SessionListResult {
	cloned := agents.CloneSessionListResult(result)
	sessionID := threadSessionID(thread.AgentOptionsJSON)
	if sessionID == "" {
		return cloned
	}

	current := agents.SessionInfo{
		SessionID: sessionID,
		CWD:       thread.CWD,
	}
	sessions := make([]agents.SessionInfo, 0, len(cloned.Sessions)+1)
	sessions = append(sessions, current)
	seen := map[string]int{sessionID: 0}
	for _, session := range cloned.Sessions {
		item := agents.CloneSessionInfo(session)
		if item.SessionID == "" {
			continue
		}
		if index, ok := seen[item.SessionID]; ok {
			sessions[index] = mergeSessionInfo(sessions[index], item)
			continue
		}
		seen[item.SessionID] = len(sessions)
		sessions = append(sessions, item)
	}
	cloned.Sessions = sessions
	return cloned
}

func toSessionInfoResponse(session agents.SessionInfo, isActive bool) sessionInfoResponse {
	cloned := agents.CloneSessionInfo(session)
	return sessionInfoResponse{
		SessionID: cloned.SessionID,
		CWD:       cloned.CWD,
		Title:     cloned.Title,
		UpdatedAt: cloned.UpdatedAt,
		Meta:      cloned.Meta,
		IsActive:  isActive,
	}
}

func (s *Server) sessionListResponseItems(thread storage.Thread, result agents.SessionListResult) []sessionInfoResponse {
	result = includeCurrentThreadSession(thread, result)
	if len(result.Sessions) == 0 {
		return []sessionInfoResponse{}
	}

	hasActiveSession := s.turns.HasActiveSession(thread.ThreadID)
	items := make([]sessionInfoResponse, 0, len(result.Sessions))
	for _, session := range result.Sessions {
		sessionID := strings.TrimSpace(session.SessionID)
		items = append(items, toSessionInfoResponse(
			session,
			hasActiveSession && s.turns.IsSessionActive(thread.ThreadID, sessionID),
		))
	}
	return items
}

func (s *Server) handleThreadSessions(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	provider, err := s.turnAgentFactory(thread)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, codeUpstreamUnavailable, "failed to resolve agent provider", map[string]any{
			"agent":  thread.AgentID,
			"reason": err.Error(),
		})
		return
	}
	if closer, ok := provider.(io.Closer); ok {
		defer closer.Close()
	}

	lister, ok := provider.(agents.SessionLister)
	if !ok {
		writeJSON(w, http.StatusOK, map[string]any{
			"threadId":   thread.ThreadID,
			"supported":  false,
			"sessions":   []agents.SessionInfo{},
			"nextCursor": "",
		})
		return
	}

	result, err := lister.ListSessions(r.Context(), agents.SessionListRequest{
		CWD:    thread.CWD,
		Cursor: strings.TrimSpace(r.URL.Query().Get("cursor")),
	})
	if err != nil {
		if errors.Is(err, agents.ErrSessionListUnsupported) || errors.Is(err, agents.ErrSessionLoadUnsupported) {
			writeJSON(w, http.StatusOK, map[string]any{
				"threadId":   thread.ThreadID,
				"supported":  false,
				"sessions":   []agents.SessionInfo{},
				"nextCursor": "",
			})
			return
		}
		writeError(w, http.StatusServiceUnavailable, codeUpstreamUnavailable, "failed to query thread sessions", map[string]any{
			"threadId": thread.ThreadID,
			"agent":    thread.AgentID,
			"reason":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"threadId":   thread.ThreadID,
		"supported":  true,
		"sessions":   s.sessionListResponseItems(thread, result),
		"nextCursor": result.NextCursor,
	})
}

func (s *Server) handleThreadSessionHistory(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	sessionID := strings.TrimSpace(r.URL.Query().Get("sessionId"))
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "sessionId is required", map[string]any{
			"field": "sessionId",
		})
		return
	}

	cachedResult, found, err := s.loadCachedSessionTranscript(r.Context(), thread, sessionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to load session transcript cache", map[string]any{
			"threadId":  thread.ThreadID,
			"sessionId": sessionID,
			"reason":    err.Error(),
		})
		return
	}
	_, configFound, configErr := s.loadStoredSessionConfigOptions(r.Context(), thread, sessionID)
	if configErr != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to load session config cache", map[string]any{
			"threadId":  thread.ThreadID,
			"sessionId": sessionID,
			"reason":    configErr.Error(),
		})
		return
	}
	if found && configFound {
		writeJSON(w, http.StatusOK, map[string]any{
			"threadId":  thread.ThreadID,
			"sessionId": sessionID,
			"supported": true,
			"messages":  cachedResult.Messages,
		})
		return
	}

	provider, err := s.turnAgentFactory(thread)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, codeUpstreamUnavailable, "failed to resolve agent provider", map[string]any{
			"agent":  thread.AgentID,
			"reason": err.Error(),
		})
		return
	}
	if closer, ok := provider.(io.Closer); ok {
		defer closer.Close()
	}

	loader, ok := provider.(agents.SessionTranscriptLoader)
	if !ok {
		if found {
			writeJSON(w, http.StatusOK, map[string]any{
				"threadId":  thread.ThreadID,
				"sessionId": sessionID,
				"supported": true,
				"messages":  cachedResult.Messages,
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"threadId":  thread.ThreadID,
			"sessionId": sessionID,
			"supported": false,
			"messages":  []agents.SessionTranscriptMessage{},
		})
		return
	}

	loadCtx := agents.WithSessionUsageHandler(r.Context(), func(sessionUsageCtx context.Context, update agents.SessionUsageUpdate) error {
		_ = sessionUsageCtx
		s.persistSessionUsageSnapshotBestEffort(r.Context(), thread, update)
		return nil
	})

	result, err := loader.LoadSessionTranscript(loadCtx, agents.SessionTranscriptRequest{
		CWD:       thread.CWD,
		SessionID: sessionID,
	})
	if err != nil {
		if found {
			writeJSON(w, http.StatusOK, map[string]any{
				"threadId":  thread.ThreadID,
				"sessionId": sessionID,
				"supported": true,
				"messages":  cachedResult.Messages,
			})
			return
		}
		if errors.Is(err, agents.ErrSessionLoadUnsupported) || errors.Is(err, agents.ErrSessionListUnsupported) {
			writeJSON(w, http.StatusOK, map[string]any{
				"threadId":  thread.ThreadID,
				"sessionId": sessionID,
				"supported": false,
				"messages":  []agents.SessionTranscriptMessage{},
			})
			return
		}
		if errors.Is(err, agents.ErrSessionNotFound) {
			writeError(w, http.StatusNotFound, codeNotFound, "session not found", map[string]any{
				"threadId":  thread.ThreadID,
				"sessionId": sessionID,
			})
			return
		}
		writeError(w, http.StatusServiceUnavailable, codeUpstreamUnavailable, "failed to load session transcript", map[string]any{
			"threadId":  thread.ThreadID,
			"sessionId": sessionID,
			"agent":     thread.AgentID,
			"reason":    err.Error(),
		})
		return
	}

	result = agents.CloneSessionTranscriptResult(result)
	s.persistSessionLoadConfigSnapshotBestEffort(r.Context(), &thread, sessionID, result.ConfigOptions)
	if err := s.persistSessionTranscriptCache(r.Context(), thread, sessionID, result); err != nil {
		s.logger.Warn("session_history.cache_persist_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"sessionId", sessionID,
			"reason", err.Error(),
		)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"threadId":  thread.ThreadID,
		"sessionId": sessionID,
		"supported": true,
		"messages":  result.Messages,
	})
}

func (s *Server) handleThreadSessionUsage(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	sessionID := strings.TrimSpace(r.URL.Query().Get("sessionId"))
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, codeInvalidArgument, "sessionId is required", map[string]any{
			"field": "sessionId",
		})
		return
	}

	cachedUsage, found, err := s.loadStoredSessionUsage(r.Context(), thread, sessionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to load session usage cache", map[string]any{
			"threadId":  thread.ThreadID,
			"sessionId": sessionID,
			"reason":    err.Error(),
		})
		return
	}

	payload := map[string]any{
		"threadId":  thread.ThreadID,
		"sessionId": sessionID,
	}
	if found {
		payload["usage"] = sessionUsageResponsePayload(cachedUsage)
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) loadCachedSessionTranscript(
	ctx context.Context,
	thread storage.Thread,
	sessionID string,
) (agents.SessionTranscriptResult, bool, error) {
	cache, err := s.store.GetSessionTranscriptCache(ctx, thread.AgentID, thread.CWD, sessionID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return agents.SessionTranscriptResult{}, false, nil
		}
		return agents.SessionTranscriptResult{}, false, err
	}

	result, err := decodeSessionTranscriptCache(cache.MessagesJSON)
	if err != nil {
		s.logger.Warn("session_history.cache_decode_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"sessionId", sessionID,
			"reason", err.Error(),
		)
		return agents.SessionTranscriptResult{}, false, nil
	}
	return result, true, nil
}

func (s *Server) persistSessionTranscriptCache(
	ctx context.Context,
	thread storage.Thread,
	sessionID string,
	result agents.SessionTranscriptResult,
) error {
	messagesJSON, err := encodeSessionTranscriptCache(result)
	if err != nil {
		return err
	}
	return s.store.UpsertSessionTranscriptCache(ctx, storage.UpsertSessionTranscriptCacheParams{
		AgentID:      thread.AgentID,
		CWD:          thread.CWD,
		SessionID:    sessionID,
		MessagesJSON: messagesJSON,
	})
}

func decodeSessionTranscriptCache(raw string) (agents.SessionTranscriptResult, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		raw = "[]"
	}

	var messages []agents.SessionTranscriptMessage
	if err := json.Unmarshal([]byte(raw), &messages); err != nil {
		return agents.SessionTranscriptResult{}, fmt.Errorf("decode session transcript cache: %w", err)
	}
	return agents.CloneSessionTranscriptResult(agents.SessionTranscriptResult{Messages: messages}), nil
}

func encodeSessionTranscriptCache(result agents.SessionTranscriptResult) (string, error) {
	result = agents.CloneSessionTranscriptResult(result)
	messages := result.Messages
	if len(messages) == 0 {
		messages = []agents.SessionTranscriptMessage{}
	}
	payload, err := json.Marshal(messages)
	if err != nil {
		return "", fmt.Errorf("encode session transcript cache: %w", err)
	}
	return string(payload), nil
}

func (s *Server) loadStoredSessionUsage(
	ctx context.Context,
	thread storage.Thread,
	sessionID string,
) (storage.SessionUsageCache, bool, error) {
	cache, err := s.store.GetSessionUsageCache(ctx, thread.AgentID, thread.CWD, sessionID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return storage.SessionUsageCache{}, false, nil
		}
		return storage.SessionUsageCache{}, false, err
	}
	return cache, true, nil
}

func (s *Server) persistSessionUsageSnapshotBestEffort(
	ctx context.Context,
	thread storage.Thread,
	update agents.SessionUsageUpdate,
) {
	update = agents.CloneSessionUsageUpdate(update)
	if update.SessionID == "" || !agents.HasSessionUsageValues(update) {
		return
	}

	if err := s.store.UpsertSessionUsageCache(ctx, storage.UpsertSessionUsageCacheParams{
		AgentID:           thread.AgentID,
		CWD:               thread.CWD,
		SessionID:         update.SessionID,
		TotalTokens:       update.TotalTokens,
		InputTokens:       update.InputTokens,
		OutputTokens:      update.OutputTokens,
		ThoughtTokens:     update.ThoughtTokens,
		CachedReadTokens:  update.CachedReadTokens,
		CachedWriteTokens: update.CachedWriteTokens,
		ContextUsed:       update.ContextUsed,
		ContextSize:       update.ContextSize,
		CostAmount:        update.CostAmount,
		CostCurrency:      update.CostCurrency,
	}); err != nil {
		s.logger.Warn("thread.session_usage_persist_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"sessionId", update.SessionID,
			"reason", err.Error(),
		)
	}
}

func sessionUsageEventPayload(turnID string, update agents.SessionUsageUpdate) map[string]any {
	update = agents.CloneSessionUsageUpdate(update)
	payload := map[string]any{
		"turnId":    strings.TrimSpace(turnID),
		"sessionId": update.SessionID,
	}
	if update.TotalTokens != nil {
		payload["totalTokens"] = *update.TotalTokens
	}
	if update.InputTokens != nil {
		payload["inputTokens"] = *update.InputTokens
	}
	if update.OutputTokens != nil {
		payload["outputTokens"] = *update.OutputTokens
	}
	if update.ThoughtTokens != nil {
		payload["thoughtTokens"] = *update.ThoughtTokens
	}
	if update.CachedReadTokens != nil {
		payload["cachedReadTokens"] = *update.CachedReadTokens
	}
	if update.CachedWriteTokens != nil {
		payload["cachedWriteTokens"] = *update.CachedWriteTokens
	}
	if update.ContextUsed != nil {
		payload["contextUsed"] = *update.ContextUsed
	}
	if update.ContextSize != nil {
		payload["contextSize"] = *update.ContextSize
	}
	if update.CostAmount != nil && update.CostCurrency != "" {
		payload["costAmount"] = *update.CostAmount
		payload["costCurrency"] = update.CostCurrency
	}
	return payload
}

func sessionUsageResponsePayload(cache storage.SessionUsageCache) map[string]any {
	payload := map[string]any{
		"sessionId": cache.SessionID,
		"updatedAt": cache.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	if cache.TotalTokens != nil {
		payload["totalTokens"] = *cache.TotalTokens
	}
	if cache.InputTokens != nil {
		payload["inputTokens"] = *cache.InputTokens
	}
	if cache.OutputTokens != nil {
		payload["outputTokens"] = *cache.OutputTokens
	}
	if cache.ThoughtTokens != nil {
		payload["thoughtTokens"] = *cache.ThoughtTokens
	}
	if cache.CachedReadTokens != nil {
		payload["cachedReadTokens"] = *cache.CachedReadTokens
	}
	if cache.CachedWriteTokens != nil {
		payload["cachedWriteTokens"] = *cache.CachedWriteTokens
	}
	if cache.ContextUsed != nil {
		payload["contextUsed"] = *cache.ContextUsed
	}
	if cache.ContextSize != nil {
		payload["contextSize"] = *cache.ContextSize
	}
	if cache.CostAmount != nil && strings.TrimSpace(cache.CostCurrency) != "" {
		payload["costAmount"] = *cache.CostAmount
		payload["costCurrency"] = strings.TrimSpace(cache.CostCurrency)
	}
	return payload
}
