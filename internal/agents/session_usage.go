package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// SessionUsageUpdate describes one cumulative ACP session-usage snapshot.
type SessionUsageUpdate struct {
	SessionID         string
	TotalTokens       *int64
	InputTokens       *int64
	OutputTokens      *int64
	ThoughtTokens     *int64
	CachedReadTokens  *int64
	CachedWriteTokens *int64
	ContextUsed       *int64
	ContextSize       *int64
	CostAmount        *float64
	CostCurrency      string
}

// SessionUsageHandler receives one cumulative ACP session-usage update.
type SessionUsageHandler func(ctx context.Context, update SessionUsageUpdate) error

type sessionUsageHandlerContextKey struct{}

// WithSessionUsageHandler binds one session-usage callback to context.
func WithSessionUsageHandler(ctx context.Context, handler SessionUsageHandler) context.Context {
	if handler == nil {
		return ctx
	}
	return context.WithValue(ctx, sessionUsageHandlerContextKey{}, handler)
}

// SessionUsageHandlerFromContext gets the session-usage callback, if present.
func SessionUsageHandlerFromContext(ctx context.Context) (SessionUsageHandler, bool) {
	if ctx == nil {
		return nil, false
	}
	handler, ok := ctx.Value(sessionUsageHandlerContextKey{}).(SessionUsageHandler)
	if !ok || handler == nil {
		return nil, false
	}
	return handler, true
}

// NotifySessionUsageUpdate reports one session-usage update to the active callback.
func NotifySessionUsageUpdate(ctx context.Context, update SessionUsageUpdate) error {
	handler, ok := SessionUsageHandlerFromContext(ctx)
	if !ok {
		return nil
	}
	update = CloneSessionUsageUpdate(update)
	if strings.TrimSpace(update.SessionID) == "" || !HasSessionUsageValues(update) {
		return nil
	}
	return handler(ctx, update)
}

// HasSessionUsageValues reports whether one session-usage payload carries any metrics.
func HasSessionUsageValues(update SessionUsageUpdate) bool {
	return update.TotalTokens != nil ||
		update.InputTokens != nil ||
		update.OutputTokens != nil ||
		update.ThoughtTokens != nil ||
		update.CachedReadTokens != nil ||
		update.CachedWriteTokens != nil ||
		update.ContextUsed != nil ||
		update.ContextSize != nil ||
		update.CostAmount != nil
}

// CloneSessionUsageUpdate returns a trimmed deep copy of one usage payload.
func CloneSessionUsageUpdate(update SessionUsageUpdate) SessionUsageUpdate {
	cloned := SessionUsageUpdate{
		SessionID:         strings.TrimSpace(update.SessionID),
		TotalTokens:       cloneInt64Ptr(update.TotalTokens),
		InputTokens:       cloneInt64Ptr(update.InputTokens),
		OutputTokens:      cloneInt64Ptr(update.OutputTokens),
		ThoughtTokens:     cloneInt64Ptr(update.ThoughtTokens),
		CachedReadTokens:  cloneInt64Ptr(update.CachedReadTokens),
		CachedWriteTokens: cloneInt64Ptr(update.CachedWriteTokens),
		ContextUsed:       cloneInt64Ptr(update.ContextUsed),
		ContextSize:       cloneInt64Ptr(update.ContextSize),
		CostAmount:        cloneFloat64Ptr(update.CostAmount),
		CostCurrency:      strings.TrimSpace(update.CostCurrency),
	}
	if cloned.CostAmount == nil || cloned.CostCurrency == "" {
		cloned.CostAmount = nil
		cloned.CostCurrency = ""
	}
	return cloned
}

// MergeSessionUsageUpdate overlays a partial patch onto a previous snapshot.
func MergeSessionUsageUpdate(base, patch SessionUsageUpdate) SessionUsageUpdate {
	base = CloneSessionUsageUpdate(base)
	patch = CloneSessionUsageUpdate(patch)

	if patch.SessionID != "" {
		base.SessionID = patch.SessionID
	}
	if patch.TotalTokens != nil {
		base.TotalTokens = cloneInt64Ptr(patch.TotalTokens)
	}
	if patch.InputTokens != nil {
		base.InputTokens = cloneInt64Ptr(patch.InputTokens)
	}
	if patch.OutputTokens != nil {
		base.OutputTokens = cloneInt64Ptr(patch.OutputTokens)
	}
	if patch.ThoughtTokens != nil {
		base.ThoughtTokens = cloneInt64Ptr(patch.ThoughtTokens)
	}
	if patch.CachedReadTokens != nil {
		base.CachedReadTokens = cloneInt64Ptr(patch.CachedReadTokens)
	}
	if patch.CachedWriteTokens != nil {
		base.CachedWriteTokens = cloneInt64Ptr(patch.CachedWriteTokens)
	}
	if patch.ContextUsed != nil {
		base.ContextUsed = cloneInt64Ptr(patch.ContextUsed)
	}
	if patch.ContextSize != nil {
		base.ContextSize = cloneInt64Ptr(patch.ContextSize)
	}
	if patch.CostAmount != nil && patch.CostCurrency != "" {
		base.CostAmount = cloneFloat64Ptr(patch.CostAmount)
		base.CostCurrency = patch.CostCurrency
	}
	return base
}

// ParseACPPromptUsage extracts one optional PromptResponse.usage payload.
func ParseACPPromptUsage(raw json.RawMessage) (SessionUsageUpdate, error) {
	var payload struct {
		SessionID string          `json:"sessionId"`
		Usage     json.RawMessage `json:"usage"`
	}
	if len(raw) == 0 {
		return SessionUsageUpdate{}, nil
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return SessionUsageUpdate{}, fmt.Errorf("decode ACP prompt response: %w", err)
	}

	update := SessionUsageUpdate{
		SessionID: strings.TrimSpace(payload.SessionID),
	}
	usageRaw := json.RawMessage(strings.TrimSpace(string(payload.Usage)))
	if len(usageRaw) == 0 || string(usageRaw) == "null" {
		return update, nil
	}

	var usage struct {
		TotalTokens       *int64 `json:"total_tokens"`
		InputTokens       *int64 `json:"input_tokens"`
		OutputTokens      *int64 `json:"output_tokens"`
		ThoughtTokens     *int64 `json:"thought_tokens"`
		CachedReadTokens  *int64 `json:"cached_read_tokens"`
		CachedWriteTokens *int64 `json:"cached_write_tokens"`
	}
	if err := json.Unmarshal(usageRaw, &usage); err != nil {
		return update, nil
	}

	update.TotalTokens = cloneInt64Ptr(usage.TotalTokens)
	update.InputTokens = cloneInt64Ptr(usage.InputTokens)
	update.OutputTokens = cloneInt64Ptr(usage.OutputTokens)
	update.ThoughtTokens = cloneInt64Ptr(usage.ThoughtTokens)
	update.CachedReadTokens = cloneInt64Ptr(usage.CachedReadTokens)
	update.CachedWriteTokens = cloneInt64Ptr(usage.CachedWriteTokens)
	return CloneSessionUsageUpdate(update), nil
}

func cloneInt64Ptr(value *int64) *int64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneFloat64Ptr(value *float64) *float64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
