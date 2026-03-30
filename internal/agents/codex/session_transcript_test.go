package codex

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/beyond5959/acp-adapter/pkg/codexacp"
	"github.com/beyond5959/ngent/internal/agents"
)

func TestConsumeCodexReplayUpdate(t *testing.T) {
	t.Parallel()

	collector := agents.NewACPTranscriptCollector()
	updates := []codexacp.RPCMessage{
		{
			Method: methodSessionUpdate,
			Params: json.RawMessage(`{"update":{"sessionUpdate":"user_message_chunk","content":{"type":"text","text":"hello codex"}}}`),
		},
		{
			Method: methodSessionUpdate,
			Params: json.RawMessage(`{"update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"reply"}}}`),
		},
	}
	for _, update := range updates {
		if err := consumeCodexReplayUpdate(context.Background(), collector, update); err != nil {
			t.Fatalf("consumeCodexReplayUpdate() error = %v", err)
		}
	}

	result := collector.Result()
	if got, want := len(result.Messages), 2; got != want {
		t.Fatalf("len(messages) = %d, want %d", got, want)
	}
	if got := result.Messages[0]; got.Role != "user" || got.Content != "hello codex" {
		t.Fatalf("messages[0] = %+v, want user hello codex", got)
	}
	if got := result.Messages[1]; got.Role != "assistant" || got.Content != "reply" {
		t.Fatalf("messages[1] = %+v, want assistant reply", got)
	}
}

func TestDrainCodexReplayUpdates(t *testing.T) {
	t.Parallel()

	collector := agents.NewACPTranscriptCollector()
	updates := make(chan codexacp.RPCMessage, 2)
	updates <- codexacp.RPCMessage{
		Method: methodSessionUpdate,
		Params: json.RawMessage(`{"update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"chunk one"}}}`),
	}
	updates <- codexacp.RPCMessage{
		Method: methodSessionUpdate,
		Params: json.RawMessage(`{"update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":" and two"}}}`),
	}

	if err := drainCodexReplayUpdates(context.Background(), collector, updates); err != nil {
		t.Fatalf("drainCodexReplayUpdates() error = %v", err)
	}

	result := collector.Result()
	if got, want := len(result.Messages), 1; got != want {
		t.Fatalf("len(messages) = %d, want %d", got, want)
	}
	if got := result.Messages[0].Content; got != "chunk one and two" {
		t.Fatalf("messages[0].Content = %q, want %q", got, "chunk one and two")
	}
}
