package opencode

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/beyond5959/ngent/internal/agents"
)

func TestHandlePermissionRequestParsesExternalDirectoryPayload(t *testing.T) {
	t.Parallel()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}
	expectedPath := filepath.Join(homeDir, ".config/opencode")

	raw := json.RawMessage(fmt.Sprintf(`{
		"sessionId": "ses_opencode_perm",
		"options": [
			{"optionId":"once","kind":"allow_once"},
			{"optionId":"always","kind":"allow_always"},
			{"optionId":"reject","kind":"reject_once"}
		],
		"toolCall": {
			"title": "external_directory",
			"kind": "other",
			"toolCallId": "call_function_test_1",
			"rawInput": {
				"filepath": %q,
				"parentDir": %q
			}
		}
	}`, expectedPath, expectedPath))

	var got agents.PermissionRequest
	resp, err := handlePermissionRequest(
		context.Background(),
		raw,
		func(ctx context.Context, req agents.PermissionRequest) (agents.PermissionResponse, error) {
			got = req
			return agents.PermissionResponse{Outcome: agents.PermissionOutcomeApproved}, nil
		},
		true,
	)
	if err != nil {
		t.Fatalf("handlePermissionRequest() error = %v", err)
	}

	if got.Approval != "file" {
		t.Fatalf("req.Approval = %q, want %q", got.Approval, "file")
	}
	if got.Command != expectedPath {
		t.Fatalf("req.Command = %q, want %q", got.Command, expectedPath)
	}
	if sessionID, _ := got.RawParams["sessionId"].(string); sessionID != "ses_opencode_perm" {
		t.Fatalf("req.RawParams[sessionId] = %q, want %q", sessionID, "ses_opencode_perm")
	}
	if path, _ := got.RawParams["path"].(string); path != expectedPath {
		t.Fatalf("req.RawParams[path] = %q, want %q", path, expectedPath)
	}

	var decoded struct {
		Outcome struct {
			Outcome  string `json:"outcome"`
			OptionID string `json:"optionId"`
		} `json:"outcome"`
	}
	if err := json.Unmarshal(resp, &decoded); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if decoded.Outcome.Outcome != "selected" {
		t.Fatalf("response.outcome.outcome = %q, want %q", decoded.Outcome.Outcome, "selected")
	}
	if decoded.Outcome.OptionID != "once" {
		t.Fatalf("response.outcome.optionId = %q, want %q", decoded.Outcome.OptionID, "once")
	}
}
