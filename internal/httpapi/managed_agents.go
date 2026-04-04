package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/storage"
)

type managedAgent struct {
	scopeKey  string
	threadID  string
	sessionID string
	provider  agents.Streamer
	closer    io.Closer
	lastUsed  time.Time
}

func normalizeThreadAgentOptionsForScope(agentOptionsJSON string) string {
	scopeOptions := map[string]any{}
	if sessionID := threadSessionID(agentOptionsJSON); sessionID != "" {
		scopeOptions["sessionId"] = sessionID
	}
	if threadFreshSessionRequested(agentOptionsJSON) {
		scopeOptions[threadAgentOptionFreshSessionKey] = true
	}
	normalized, err := json.Marshal(scopeOptions)
	if err != nil {
		return "{}"
	}
	return string(normalized)
}

func threadAgentScopeKey(thread storage.Thread) string {
	return threadAgentScopeKeyFromOptions(thread.ThreadID, thread.AgentOptionsJSON)
}

func threadAgentScopeKeyFromOptions(threadID, agentOptionsJSON string) string {
	return threadID + "\x00" + normalizeThreadAgentOptionsForScope(agentOptionsJSON)
}

func (s *Server) resolveTurnAgent(thread storage.Thread) (agents.Streamer, error) {
	scopeKey := threadAgentScopeKey(thread)
	sessionID := threadSessionID(thread.AgentOptionsJSON)
	s.agentMu.Lock()
	entry, ok := s.agentsByScope[scopeKey]
	if ok {
		entry.lastUsed = time.Now().UTC()
		provider := entry.provider
		s.agentMu.Unlock()
		return provider, nil
	}
	s.agentMu.Unlock()

	if s.turnAgentFactory == nil {
		return nil, errors.New("turn agent factory is not configured")
	}
	provider, err := s.turnAgentFactory(thread)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, errors.New("turn agent factory returned nil provider")
	}

	var closer io.Closer
	if c, ok := provider.(io.Closer); ok {
		closer = c
	}

	s.agentMu.Lock()
	if existing, exists := s.agentsByScope[scopeKey]; exists {
		existing.lastUsed = time.Now().UTC()
		s.agentMu.Unlock()
		if closer != nil {
			_ = closer.Close()
		}
		return existing.provider, nil
	}
	s.agentsByScope[scopeKey] = &managedAgent{
		scopeKey:  scopeKey,
		threadID:  thread.ThreadID,
		sessionID: sessionID,
		provider:  provider,
		closer:    closer,
		lastUsed:  time.Now().UTC(),
	}
	s.agentMu.Unlock()
	return provider, nil
}

func (s *Server) Close() error {
	select {
	case <-s.janitorStop:
	default:
		close(s.janitorStop)
	}
	<-s.janitorDone
	return s.closeAllThreadAgents()
}

func (s *Server) idleJanitorLoop() {
	defer close(s.janitorDone)
	interval := s.agentIdleTTL / 2
	if interval < 500*time.Millisecond {
		interval = 500 * time.Millisecond
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.janitorStop:
			return
		case <-ticker.C:
			s.reapIdleAgents(time.Now().UTC())
		}
	}
}

func (s *Server) reapIdleAgents(now time.Time) {
	if s.agentIdleTTL <= 0 {
		return
	}

	type reclaimItem struct {
		threadID  string
		scopeKey  string
		sessionID string
		name      string
		idleFor   time.Duration
		closer    io.Closer
	}
	items := make([]reclaimItem, 0)

	s.agentMu.Lock()
	for scopeKey, entry := range s.agentsByScope {
		if s.turns.IsSessionActive(entry.threadID, entry.sessionID) {
			continue
		}
		idleFor := now.Sub(entry.lastUsed)
		if idleFor < s.agentIdleTTL {
			continue
		}
		delete(s.agentsByScope, scopeKey)
		items = append(items, reclaimItem{
			threadID:  entry.threadID,
			scopeKey:  scopeKey,
			sessionID: entry.sessionID,
			name:      entry.provider.Name(),
			idleFor:   idleFor,
			closer:    entry.closer,
		})
	}
	s.agentMu.Unlock()

	for _, item := range items {
		if item.closer != nil {
			_ = item.closer.Close()
		}
		s.logger.Info("agent.idle_reclaimed",
			"threadId", item.threadID,
			"sessionId", item.sessionID,
			"agentName", item.name,
			"idleFor", item.idleFor.String(),
		)
	}
}

func (s *Server) closeAllThreadAgents() error {
	type closeItem struct {
		threadID  string
		sessionID string
		name      string
		closer    io.Closer
	}

	items := make([]closeItem, 0)
	s.agentMu.Lock()
	for scopeKey, entry := range s.agentsByScope {
		items = append(items, closeItem{
			threadID:  entry.threadID,
			sessionID: entry.sessionID,
			name:      entry.provider.Name(),
			closer:    entry.closer,
		})
		delete(s.agentsByScope, scopeKey)
	}
	s.agentMu.Unlock()

	for _, item := range items {
		if item.closer != nil {
			_ = item.closer.Close()
		}
		s.logger.Info("agent.closed",
			"threadId", item.threadID,
			"sessionId", item.sessionID,
			"agentName", item.name,
			"reason", "server_close",
		)
	}
	return nil
}

func (s *Server) closeThreadAgents(threadID, reason string) {
	if strings.TrimSpace(threadID) == "" {
		return
	}

	type closeItem struct {
		sessionID string
		name      string
		closer    io.Closer
	}

	items := make([]closeItem, 0)
	s.agentMu.Lock()
	for scopeKey, entry := range s.agentsByScope {
		if entry.threadID != threadID {
			continue
		}
		items = append(items, closeItem{
			sessionID: entry.sessionID,
			name:      entry.provider.Name(),
			closer:    entry.closer,
		})
		delete(s.agentsByScope, scopeKey)
	}
	s.agentMu.Unlock()

	for _, item := range items {
		if item.closer != nil {
			_ = item.closer.Close()
		}
		s.logger.Info("agent.closed",
			"threadId", threadID,
			"sessionId", item.sessionID,
			"agentName", item.name,
			"reason", reason,
		)
	}
}

func (s *Server) closeThreadAgentScope(threadID, agentOptionsJSON, reason string) {
	threadID = strings.TrimSpace(threadID)
	if threadID == "" {
		return
	}

	sessionID := threadSessionID(agentOptionsJSON)
	if s.turns.IsSessionActive(threadID, sessionID) {
		return
	}

	scopeKey := threadAgentScopeKeyFromOptions(threadID, agentOptionsJSON)

	var item *managedAgent
	s.agentMu.Lock()
	entry, ok := s.agentsByScope[scopeKey]
	if ok {
		delete(s.agentsByScope, scopeKey)
		item = entry
	}
	s.agentMu.Unlock()
	if item == nil {
		return
	}

	if item.closer != nil {
		_ = item.closer.Close()
	}
	s.logger.Info("agent.closed",
		"threadId", item.threadID,
		"sessionId", item.sessionID,
		"agentName", item.provider.Name(),
		"reason", reason,
	)
}

func (s *Server) rebindManagedAgentScope(threadID, fromAgentOptionsJSON, toAgentOptionsJSON string) error {
	threadID = strings.TrimSpace(threadID)
	if threadID == "" {
		return nil
	}
	fromScopeKey := threadAgentScopeKeyFromOptions(threadID, fromAgentOptionsJSON)
	toScopeKey := threadAgentScopeKeyFromOptions(threadID, toAgentOptionsJSON)
	if fromScopeKey == toScopeKey {
		return nil
	}

	s.agentMu.Lock()
	defer s.agentMu.Unlock()

	entry, ok := s.agentsByScope[fromScopeKey]
	if !ok {
		return nil
	}
	if _, exists := s.agentsByScope[toScopeKey]; exists {
		return nil
	}

	delete(s.agentsByScope, fromScopeKey)
	entry.scopeKey = toScopeKey
	entry.sessionID = threadSessionID(toAgentOptionsJSON)
	entry.lastUsed = time.Now().UTC()
	s.agentsByScope[toScopeKey] = entry
	return nil
}
