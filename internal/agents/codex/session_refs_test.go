package codex

import (
	"context"
	"testing"

	"github.com/beyond5959/ngent/internal/agents"
)

func TestNormalizeCodexSessionListResultUsesStableThreadID(t *testing.T) {
	result := normalizeCodexSessionListResult(agents.SessionListResult{
		Sessions: []agents.SessionInfo{
			{
				SessionID: "session-1",
				Meta: map[string]any{
					codexMetaThreadID: "thread-123",
				},
			},
		},
	})

	if got, want := len(result.Sessions), 1; got != want {
		t.Fatalf("len(sessions) = %d, want %d", got, want)
	}
	if got, want := result.Sessions[0].SessionID, "thread-123"; got != want {
		t.Fatalf("sessionId = %q, want %q", got, want)
	}
	if got, want := codexLoadSessionID(result.Sessions[0]), "session-1"; got != want {
		t.Fatalf("load session id = %q, want %q", got, want)
	}
}

func TestCodexSessionMatchesIDAcceptsStableAndRawIDs(t *testing.T) {
	session := normalizeCodexSessionInfo(agents.SessionInfo{
		SessionID: "session-7",
		Meta: map[string]any{
			codexMetaThreadID: "thread-789",
		},
	})

	if !codexSessionMatchesID(session, "thread-789") {
		t.Fatalf("stable session id did not match")
	}
	if !codexSessionMatchesID(session, "session-7") {
		t.Fatalf("raw session id did not match")
	}
	if codexSessionMatchesID(session, "session-8") {
		t.Fatalf("unexpected match for unrelated session id")
	}
}

func TestCodexStableSessionIDFallsBackToRawSessionID(t *testing.T) {
	session := normalizeCodexSessionInfo(agents.SessionInfo{SessionID: "session-9"})

	if got, want := session.SessionID, "session-9"; got != want {
		t.Fatalf("sessionId = %q, want %q", got, want)
	}
	if got, want := codexLoadSessionID(session), "session-9"; got != want {
		t.Fatalf("load session id = %q, want %q", got, want)
	}
}

func TestCodexShouldDeferInitialSessionBinding(t *testing.T) {
	if !codexShouldDeferInitialSessionBinding("", "session-1", "session-1") {
		t.Fatalf("expected provisional new-session binding to defer")
	}
	if codexShouldDeferInitialSessionBinding("thread-123", "session-1", "session-1") {
		t.Fatalf("did not expect loaded session binding to defer")
	}
	if codexShouldDeferInitialSessionBinding("", "session-1", "thread-123") {
		t.Fatalf("did not expect stable session binding to defer")
	}
}

func TestNormalizeCodexSessionUsageUpdateUsesStableID(t *testing.T) {
	update := normalizeCodexSessionUsageUpdate(agents.SessionUsageUpdate{
		SessionID:   "session-1",
		TotalTokens: int64Ptr(42),
	}, "thread-123", "session-1")

	if got, want := update.SessionID, "thread-123"; got != want {
		t.Fatalf("sessionId = %q, want %q", got, want)
	}
	if got, want := *update.TotalTokens, int64(42); got != want {
		t.Fatalf("totalTokens = %d, want %d", got, want)
	}
}

func TestNotifyCachedSessionUsagePromotesRawID(t *testing.T) {
	client := &Client{
		sessionUsageByID: map[string]agents.SessionUsageUpdate{
			"session-1": {
				SessionID:   "session-1",
				ContextUsed: int64Ptr(53000),
				ContextSize: int64Ptr(200000),
			},
		},
	}

	var captured agents.SessionUsageUpdate
	ctx := agents.WithSessionUsageHandler(context.Background(), func(_ context.Context, update agents.SessionUsageUpdate) error {
		captured = agents.CloneSessionUsageUpdate(update)
		return nil
	})

	if err := client.notifyCachedSessionUsage(ctx, "thread-123", "session-1"); err != nil {
		t.Fatalf("notifyCachedSessionUsage() error = %v", err)
	}

	if got, want := captured.SessionID, "thread-123"; got != want {
		t.Fatalf("reported sessionId = %q, want %q", got, want)
	}
	if got, want := *captured.ContextUsed, int64(53000); got != want {
		t.Fatalf("contextUsed = %d, want %d", got, want)
	}
	if got, want := *captured.ContextSize, int64(200000); got != want {
		t.Fatalf("contextSize = %d, want %d", got, want)
	}
	if _, ok := client.sessionUsageByID["session-1"]; ok {
		t.Fatalf("raw session usage cache entry still present")
	}
	if got := client.sessionUsageByID["thread-123"].SessionID; got != "thread-123" {
		t.Fatalf("stable cache entry sessionId = %q, want %q", got, "thread-123")
	}
}

func int64Ptr(value int64) *int64 {
	return &value
}
