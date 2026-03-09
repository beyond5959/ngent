package agents

import (
	"encoding/json"
	"testing"
)

func TestParseACPUpdateMessageChunk(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`{
		"update": {
			"sessionUpdate": "agent_message_chunk",
			"content": {
				"type": "text",
				"text": "hello"
			}
		}
	}`)

	update, err := ParseACPUpdate(raw)
	if err != nil {
		t.Fatalf("ParseACPUpdate() error = %v", err)
	}
	if update.Type != ACPUpdateTypeMessageChunk {
		t.Fatalf("update.Type = %q, want %q", update.Type, ACPUpdateTypeMessageChunk)
	}
	if update.Delta != "hello" {
		t.Fatalf("update.Delta = %q, want %q", update.Delta, "hello")
	}
	if len(update.PlanEntries) != 0 {
		t.Fatalf("len(update.PlanEntries) = %d, want 0", len(update.PlanEntries))
	}
}

func TestParseACPUpdatePlan(t *testing.T) {
	t.Parallel()

	raw := json.RawMessage(`{
		"update": {
			"sessionUpdate": "plan",
			"entries": [
				{"content": " Inspect files ", "status": " in_progress ", "priority": " high "},
				{"content": "   ", "status": "pending", "priority": "low"},
				{"content": "Write patch", "status": "pending", "priority": "medium"}
			]
		}
	}`)

	update, err := ParseACPUpdate(raw)
	if err != nil {
		t.Fatalf("ParseACPUpdate() error = %v", err)
	}
	if update.Type != ACPUpdateTypePlan {
		t.Fatalf("update.Type = %q, want %q", update.Type, ACPUpdateTypePlan)
	}
	if len(update.PlanEntries) != 2 {
		t.Fatalf("len(update.PlanEntries) = %d, want 2", len(update.PlanEntries))
	}
	if got := update.PlanEntries[0]; got.Content != "Inspect files" || got.Status != "in_progress" || got.Priority != "high" {
		t.Fatalf("first plan entry = %+v, want trimmed values", got)
	}
}

func TestParseACPUpdateLegacyDelta(t *testing.T) {
	t.Parallel()

	update, err := ParseACPUpdate(json.RawMessage(`{"delta":"legacy"}`))
	if err != nil {
		t.Fatalf("ParseACPUpdate() error = %v", err)
	}
	if update.Type != ACPUpdateTypeMessageChunk {
		t.Fatalf("update.Type = %q, want %q", update.Type, ACPUpdateTypeMessageChunk)
	}
	if update.Delta != "legacy" {
		t.Fatalf("update.Delta = %q, want %q", update.Delta, "legacy")
	}
}

func TestParseACPUpdatePlanKeepsEmptyReplacement(t *testing.T) {
	t.Parallel()

	update, err := ParseACPUpdate(json.RawMessage(`{"update":{"sessionUpdate":"plan","entries":[]}}`))
	if err != nil {
		t.Fatalf("ParseACPUpdate() error = %v", err)
	}
	if update.Type != ACPUpdateTypePlan {
		t.Fatalf("update.Type = %q, want %q", update.Type, ACPUpdateTypePlan)
	}
	if update.PlanEntries != nil {
		t.Fatalf("update.PlanEntries = %#v, want nil empty replacement", update.PlanEntries)
	}
}
