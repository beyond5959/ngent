package runtime

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	// ErrActiveTurnExists means another turn is currently running on the scope.
	ErrActiveTurnExists = errors.New("runtime: active turn already exists for scope")
	// ErrTurnNotActive means the turn is not tracked as active.
	ErrTurnNotActive = errors.New("runtime: turn is not active")
)

type activeTurn struct {
	threadID        string
	sessionID       string
	scopeKey        string
	turnID          string
	cancel          context.CancelFunc
	threadExclusive bool
}

// TurnController manages active turn lifecycle and cancellation.
type TurnController struct {
	mu           sync.Mutex
	cond         *sync.Cond
	byScope      map[string]activeTurn
	byTurn       map[string]activeTurn
	threadActive map[string]int
	threadGuards map[string]activeTurn
}

// NewTurnController constructs a new active-turn controller.
func NewTurnController() *TurnController {
	controller := &TurnController{
		byScope:      make(map[string]activeTurn),
		byTurn:       make(map[string]activeTurn),
		threadActive: make(map[string]int),
		threadGuards: make(map[string]activeTurn),
	}
	controller.cond = sync.NewCond(&controller.mu)
	return controller
}

func turnScopeKey(threadID, sessionID string) string {
	return threadID + "\x00" + strings.TrimSpace(sessionID)
}

// Activate registers a running turn; one active turn is allowed per thread/session scope.
func (c *TurnController) Activate(threadID, sessionID, turnID string, cancel context.CancelFunc) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.threadGuards[threadID]; exists {
		return ErrActiveTurnExists
	}
	scopeKey := turnScopeKey(threadID, sessionID)
	if _, exists := c.byScope[scopeKey]; exists {
		return ErrActiveTurnExists
	}

	entry := activeTurn{
		threadID:  threadID,
		sessionID: strings.TrimSpace(sessionID),
		scopeKey:  scopeKey,
		turnID:    turnID,
		cancel:    cancel,
	}
	c.byScope[scopeKey] = entry
	c.byTurn[turnID] = entry
	c.threadActive[threadID]++
	return nil
}

// ActivateThreadExclusive blocks all scopes on the thread until ReleaseThreadExclusive.
func (c *TurnController) ActivateThreadExclusive(threadID, turnID string, cancel context.CancelFunc) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.threadGuards[threadID]; exists {
		return ErrActiveTurnExists
	}
	if c.threadActive[threadID] > 0 {
		return ErrActiveTurnExists
	}

	entry := activeTurn{
		threadID:        threadID,
		turnID:          turnID,
		cancel:          cancel,
		threadExclusive: true,
	}
	c.threadGuards[threadID] = entry
	c.byTurn[turnID] = entry
	return nil
}

// BindTurnSession updates the active session scope for one running turn.
func (c *TurnController) BindTurnSession(turnID, sessionID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.byTurn[turnID]
	if !ok {
		return ErrTurnNotActive
	}
	if entry.threadExclusive {
		return nil
	}

	nextScopeKey := turnScopeKey(entry.threadID, sessionID)
	if nextScopeKey == entry.scopeKey {
		return nil
	}
	if _, exists := c.byScope[nextScopeKey]; exists {
		return ErrActiveTurnExists
	}

	delete(c.byScope, entry.scopeKey)
	entry.sessionID = strings.TrimSpace(sessionID)
	entry.scopeKey = nextScopeKey
	c.byScope[nextScopeKey] = entry
	c.byTurn[turnID] = entry
	return nil
}

// Release removes the running turn from controller maps.
func (c *TurnController) Release(threadID, sessionID, turnID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.byTurn[turnID]
	if !ok || entry.threadExclusive {
		return
	}
	if entry.threadID != threadID {
		return
	}
	if entry.scopeKey != turnScopeKey(threadID, sessionID) {
		return
	}

	delete(c.byTurn, turnID)
	delete(c.byScope, entry.scopeKey)
	if remaining := c.threadActive[threadID] - 1; remaining > 0 {
		c.threadActive[threadID] = remaining
	} else {
		delete(c.threadActive, threadID)
	}
	c.cond.Broadcast()
}

// ReleaseThreadExclusive removes the thread-wide exclusive guard.
func (c *TurnController) ReleaseThreadExclusive(threadID, turnID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.byTurn[turnID]
	if !ok || !entry.threadExclusive {
		return
	}
	if entry.threadID != threadID {
		return
	}

	delete(c.byTurn, turnID)
	delete(c.threadGuards, threadID)
	c.cond.Broadcast()
}

// Cancel requests cancellation for an active turn.
func (c *TurnController) Cancel(turnID string) error {
	c.mu.Lock()
	entry, ok := c.byTurn[turnID]
	c.mu.Unlock()
	if !ok {
		return ErrTurnNotActive
	}

	if entry.cancel != nil {
		entry.cancel()
	}
	return nil
}

// IsThreadActive reports whether a thread has an active turn.
func (c *TurnController) IsThreadActive(threadID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.threadGuards[threadID]; ok {
		return true
	}
	return c.threadActive[threadID] > 0
}

// HasActiveSession reports whether a thread has any active session-scoped turn.
func (c *TurnController) HasActiveSession(threadID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.threadActive[threadID] > 0
}

// IsSessionActive reports whether one thread/session scope has an active turn.
func (c *TurnController) IsSessionActive(threadID, sessionID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.threadGuards[threadID]; ok {
		return true
	}
	_, ok := c.byScope[turnScopeKey(threadID, sessionID)]
	return ok
}

// ActiveCount returns currently active turn count.
func (c *TurnController) ActiveCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.byTurn)
}

// IsTurnActive reports whether one tracked turn is still active.
func (c *TurnController) IsTurnActive(turnID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.byTurn[turnID]
	return ok
}

// CancelAll requests cancellation for all active turns.
func (c *TurnController) CancelAll() int {
	c.mu.Lock()
	entries := make([]activeTurn, 0, len(c.byTurn))
	for _, entry := range c.byTurn {
		entries = append(entries, entry)
	}
	c.mu.Unlock()

	cancelled := 0
	for _, entry := range entries {
		if entry.cancel != nil {
			entry.cancel()
			cancelled++
		}
	}
	return cancelled
}

// WaitForIdle blocks until no active turns remain or context is cancelled.
func (c *TurnController) WaitForIdle(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for {
		c.mu.Lock()
		idle := len(c.byTurn) == 0
		c.mu.Unlock()
		if idle {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
