package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	// ACPUpdateTypeMessageChunk streams agent text deltas.
	ACPUpdateTypeMessageChunk = "agent_message_chunk"
	// ACPUpdateTypePlan replaces the current agent plan entries.
	ACPUpdateTypePlan = "plan"
)

// PlanEntry is one ACP plan entry shown to the user.
type PlanEntry struct {
	Content  string `json:"content"`
	Status   string `json:"status,omitempty"`
	Priority string `json:"priority,omitempty"`
}

// ACPUpdate is one normalized ACP session/update payload.
type ACPUpdate struct {
	Type        string
	Delta       string
	PlanEntries []PlanEntry
}

// ParseACPUpdate normalizes provider-specific session/update payloads.
func ParseACPUpdate(raw json.RawMessage) (ACPUpdate, error) {
	if len(raw) == 0 {
		return ACPUpdate{}, nil
	}

	var payload struct {
		Delta  string `json:"delta"`
		Update struct {
			SessionUpdate string `json:"sessionUpdate"`
			Content       struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
			Entries []PlanEntry `json:"entries"`
		} `json:"update"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ACPUpdate{}, fmt.Errorf("decode ACP session/update payload: %w", err)
	}

	switch strings.TrimSpace(payload.Update.SessionUpdate) {
	case "":
		if payload.Delta == "" {
			return ACPUpdate{}, nil
		}
		return ACPUpdate{
			Type:  ACPUpdateTypeMessageChunk,
			Delta: payload.Delta,
		}, nil
	case ACPUpdateTypeMessageChunk:
		if contentType := strings.TrimSpace(payload.Update.Content.Type); contentType != "" && contentType != "text" {
			return ACPUpdate{Type: ACPUpdateTypeMessageChunk}, nil
		}
		return ACPUpdate{
			Type:  ACPUpdateTypeMessageChunk,
			Delta: payload.Update.Content.Text,
		}, nil
	case ACPUpdateTypePlan:
		return ACPUpdate{
			Type:        ACPUpdateTypePlan,
			PlanEntries: normalizePlanEntries(payload.Update.Entries),
		}, nil
	default:
		return ACPUpdate{Type: strings.TrimSpace(payload.Update.SessionUpdate)}, nil
	}
}

// ClonePlanEntries returns a trimmed deep copy of the provided entries.
func ClonePlanEntries(entries []PlanEntry) []PlanEntry {
	return normalizePlanEntries(entries)
}

func normalizePlanEntries(entries []PlanEntry) []PlanEntry {
	if len(entries) == 0 {
		return nil
	}

	normalized := make([]PlanEntry, 0, len(entries))
	for _, entry := range entries {
		content := strings.TrimSpace(entry.Content)
		if content == "" {
			continue
		}
		normalized = append(normalized, PlanEntry{
			Content:  content,
			Status:   strings.TrimSpace(entry.Status),
			Priority: strings.TrimSpace(entry.Priority),
		})
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}

// PlanHandler receives ACP plan replacements for the active turn.
type PlanHandler func(ctx context.Context, entries []PlanEntry) error

type planHandlerContextKey struct{}

// WithPlanHandler binds one per-turn plan callback to context.
func WithPlanHandler(ctx context.Context, handler PlanHandler) context.Context {
	if handler == nil {
		return ctx
	}
	return context.WithValue(ctx, planHandlerContextKey{}, handler)
}

// PlanHandlerFromContext gets plan callback from context, if present.
func PlanHandlerFromContext(ctx context.Context) (PlanHandler, bool) {
	if ctx == nil {
		return nil, false
	}
	handler, ok := ctx.Value(planHandlerContextKey{}).(PlanHandler)
	if !ok || handler == nil {
		return nil, false
	}
	return handler, true
}
