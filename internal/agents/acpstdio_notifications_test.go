package agents

import (
	"context"
	"encoding/json"
	"testing"
)

func TestNewACPNotificationHandlerRoutesThoughtChunksToReasoningHandler(t *testing.T) {
	t.Parallel()

	var answer string
	var reasoning string
	ctx := WithReasoningHandler(context.Background(), func(ctx context.Context, delta string) error {
		_ = ctx
		reasoning += delta
		return nil
	})

	handler, markPromptStarted := NewACPNotificationHandler(ctx, func(delta string) error {
		answer += delta
		return nil
	})
	markPromptStarted()

	raw := json.RawMessage(`{
		"update": {
			"sessionUpdate": "agent_thought_chunk",
			"content": {
				"type": "text",
				"text": "thinking"
			}
		}
	}`)
	if err := handler("session/update", raw); err != nil {
		t.Fatalf("handler() error = %v", err)
	}

	if answer != "" {
		t.Fatalf("answer = %q, want empty", answer)
	}
	if reasoning != "thinking" {
		t.Fatalf("reasoning = %q, want %q", reasoning, "thinking")
	}
}
