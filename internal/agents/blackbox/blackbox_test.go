package blackbox_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/beyond5959/ngent/internal/agents"
	blackbox "github.com/beyond5959/ngent/internal/agents/blackbox"
)

// TestPreflight verifies that Preflight returns nil when the blackbox binary exists.
func TestPreflight(t *testing.T) {
	if _, err := exec.LookPath("blackbox"); err != nil {
		t.Skip("blackbox not in PATH")
	}
	if err := blackbox.Preflight(); err != nil {
		t.Fatalf("Preflight() = %v, want nil", err)
	}
}

func TestStreamWithFakeProcess(t *testing.T) {
	python3, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not in PATH")
	}

	fakeScript := fmt.Sprintf(`#!%s
import json
import sys

def noise():
    print("Amplitude Logger [Error]: Unexpected error occurred", flush=True)

def send(obj):
    sys.stdout.write(json.dumps(obj) + "\n")
    sys.stdout.flush()

noise()
for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    req = json.loads(line)
    method = req.get("method", "")
    rid = req.get("id")
    params = req.get("params", {})

    if method == "initialize":
        noise()
        send({"jsonrpc":"2.0","id":rid,"result":{
            "protocolVersion":1,
            "authMethods":[],
            "agentCapabilities":{"loadSession":False}
        }})
    elif method == "session/new":
        noise()
        send({"jsonrpc":"2.0","id":rid,"result":{
            "sessionId":"ses_blackbox_test_123"
        }})
    elif method == "session/prompt":
        sid = params.get("sessionId","")
        noise()
        send({"jsonrpc":"2.0","method":"session/update","params":{
            "sessionId":sid,
            "update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"Hello"}}
        }})
        send({"jsonrpc":"2.0","id":rid,"result":{"stopReason":"end_turn"}})
        sys.exit(0)
    elif method == "session/cancel":
        send({"jsonrpc":"2.0","id":rid,"result":{}})
        sys.exit(0)
`, python3)

	tmpDir := t.TempDir()
	fakeBin := tmpDir + "/blackbox"
	if err := os.WriteFile(fakeBin, []byte(fakeScript), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	t.Setenv("PATH", tmpDir+":"+os.Getenv("PATH"))

	c, err := blackbox.New(blackbox.Config{Dir: tmpDir})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	var deltas []string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reason, err := c.Stream(ctx, "say hello", func(delta string) error {
		deltas = append(deltas, delta)
		return nil
	})
	if err != nil {
		t.Fatalf("Stream: %v", err)
	}
	if reason != "end_turn" {
		t.Fatalf("StopReason = %q, want %q", reason, "end_turn")
	}
	if got := strings.Join(deltas, ""); !strings.Contains(got, "Hello") {
		t.Fatalf("deltas = %q, want to contain %q", got, "Hello")
	}
}

func TestStreamWithFakeProcessModelID(t *testing.T) {
	python3, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not in PATH")
	}

	fakeScript := fmt.Sprintf(`#!%s
import json
import sys

expected_model = "blackbox/test-model"

def arg_value(flag):
    try:
        idx = sys.argv.index(flag)
    except ValueError:
        return ""
    if idx + 1 >= len(sys.argv):
        return ""
    return sys.argv[idx + 1]

def send(obj):
    sys.stdout.write(json.dumps(obj) + "\n")
    sys.stdout.flush()

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    req = json.loads(line)
    method = req.get("method", "")
    rid = req.get("id")
    params = req.get("params", {})

    if method == "initialize":
        send({"jsonrpc":"2.0","id":rid,"result":{
            "protocolVersion":1,
            "authMethods":[],
            "agentCapabilities":{"loadSession":False}
        }})
    elif method == "session/new":
        if arg_value("--model") != expected_model:
            send({"jsonrpc":"2.0","id":rid,"error":{"code":-32000,"message":"startup model flag missing"}})
            sys.exit(0)
        if params.get("model","") != expected_model:
            send({"jsonrpc":"2.0","id":rid,"error":{"code":-32000,"message":"session/new model missing"}})
            sys.exit(0)
        if params.get("modelId","") != expected_model:
            send({"jsonrpc":"2.0","id":rid,"error":{"code":-32000,"message":"session/new modelId missing"}})
            sys.exit(0)
        send({"jsonrpc":"2.0","id":rid,"result":{
            "sessionId":"ses_blackbox_model_123"
        }})
    elif method == "session/prompt":
        if params.get("model","") != expected_model:
            send({"jsonrpc":"2.0","id":rid,"error":{"code":-32000,"message":"session/prompt model missing"}})
            sys.exit(0)
        sid = params.get("sessionId","")
        send({"jsonrpc":"2.0","method":"session/update","params":{
            "sessionId":sid,
            "update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"MODEL_OK"}}
        }})
        send({"jsonrpc":"2.0","id":rid,"result":{"stopReason":"end_turn"}})
        sys.exit(0)
    elif method == "session/cancel":
        send({"jsonrpc":"2.0","id":rid,"result":{}})
        sys.exit(0)
`, python3)

	tmpDir := t.TempDir()
	fakeBin := tmpDir + "/blackbox"
	if err := os.WriteFile(fakeBin, []byte(fakeScript), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	t.Setenv("PATH", tmpDir+":"+os.Getenv("PATH"))

	c, err := blackbox.New(blackbox.Config{
		Dir:     tmpDir,
		ModelID: "blackbox/test-model",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	var deltas []string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reason, err := c.Stream(ctx, "say MODEL_OK", func(delta string) error {
		deltas = append(deltas, delta)
		return nil
	})
	if err != nil {
		t.Fatalf("Stream: %v", err)
	}
	if reason != "end_turn" {
		t.Fatalf("StopReason = %q, want %q", reason, "end_turn")
	}
	if got := strings.Join(deltas, ""); !strings.Contains(got, "MODEL_OK") {
		t.Fatalf("deltas = %q, want to contain %q", got, "MODEL_OK")
	}
}

func TestListSessionsUnsupportedWithoutLoadCapability(t *testing.T) {
	python3, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not in PATH")
	}

	fakeScript := fmt.Sprintf(`#!%s
import json
import sys

def send(obj):
    sys.stdout.write(json.dumps(obj) + "\n")
    sys.stdout.flush()

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    req = json.loads(line)
    if req.get("method", "") == "initialize":
        send({"jsonrpc":"2.0","id":req.get("id"),"result":{
            "protocolVersion":1,
            "authMethods":[],
            "agentCapabilities":{"loadSession":False}
        }})
`, python3)

	tmpDir := t.TempDir()
	fakeBin := tmpDir + "/blackbox"
	if err := os.WriteFile(fakeBin, []byte(fakeScript), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	t.Setenv("PATH", tmpDir+":"+os.Getenv("PATH"))

	client, err := blackbox.New(blackbox.Config{Dir: tmpDir})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = client.ListSessions(context.Background(), agents.SessionListRequest{CWD: tmpDir})
	if !errors.Is(err, agents.ErrSessionListUnsupported) {
		t.Fatalf("ListSessions error = %v, want %v", err, agents.ErrSessionListUnsupported)
	}
}

func TestPermissionRequestsUseStructuredACPFlow(t *testing.T) {
	python3, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not in PATH")
	}

	fakeScript := fmt.Sprintf(`#!%s
import json
import sys

def send(obj):
    sys.stdout.write(json.dumps(obj) + "\n")
    sys.stdout.flush()

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    req = json.loads(line)
    method = req.get("method", "")
    rid = req.get("id")
    params = req.get("params", {})

    if method == "initialize":
        send({"jsonrpc":"2.0","id":rid,"result":{
            "protocolVersion":1,
            "authMethods":[],
            "agentCapabilities":{"loadSession":False}
        }})
    elif method == "session/new":
        send({"jsonrpc":"2.0","id":rid,"result":{
            "sessionId":"ses_blackbox_perm_123"
        }})
    elif method == "session/prompt":
        sid = params.get("sessionId","")
        perm_id = 99
        send({"jsonrpc":"2.0","id":perm_id,"method":"session/request_permission","params":{
            "sessionId":sid,
            "toolCall":{
                "title":"Execute command",
                "kind":"command",
                "toolCallId":"tool_blackbox_perm_1",
                "content":[{"type":"command","command":"echo hello"}]
            },
            "options":[
                {"optionId":"allow_once_opt","kind":"allow_once"},
                {"optionId":"reject_once_opt","kind":"reject_once"}
            ]
        }})
        marker = ""
        for rline in sys.stdin:
            rline = rline.strip()
            if not rline:
                continue
            resp = json.loads(rline)
            if resp.get("id") != perm_id:
                continue
            result = resp.get("result", {})
            outcome = result.get("outcome", {})
            if outcome.get("outcome") == "selected":
                marker = outcome.get("optionId", "")
            else:
                marker = outcome.get("outcome", "")
            break

        send({"jsonrpc":"2.0","method":"session/update","params":{
            "sessionId":sid,
            "update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":marker}}
        }})
        send({"jsonrpc":"2.0","id":rid,"result":{"stopReason":"end_turn"}})
        sys.exit(0)
    elif method == "session/cancel":
        send({"jsonrpc":"2.0","id":rid,"result":{}})
        sys.exit(0)
`, python3)

	tmpDir := t.TempDir()
	fakeBin := tmpDir + "/blackbox"
	if err := os.WriteFile(fakeBin, []byte(fakeScript), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	t.Setenv("PATH", tmpDir+":"+os.Getenv("PATH"))

	tests := []struct {
		name       string
		outcome    agents.PermissionOutcome
		wantMarker string
	}{
		{
			name:       "approved maps to allow_once",
			outcome:    agents.PermissionOutcomeApproved,
			wantMarker: "allow_once_opt",
		},
		{
			name:       "declined maps to reject_once",
			outcome:    agents.PermissionOutcomeDeclined,
			wantMarker: "reject_once_opt",
		},
		{
			name:       "cancelled maps to cancelled outcome",
			outcome:    agents.PermissionOutcomeCancelled,
			wantMarker: "cancelled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := blackbox.New(blackbox.Config{Dir: tmpDir})
			if err != nil {
				t.Fatalf("New: %v", err)
			}

			baseCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			ctx := agents.WithPermissionHandler(baseCtx, func(ctx context.Context, req agents.PermissionRequest) (agents.PermissionResponse, error) {
				return agents.PermissionResponse{Outcome: tt.outcome}, nil
			})

			var deltas []string
			reason, err := client.Stream(ctx, "permission test", func(delta string) error {
				deltas = append(deltas, delta)
				return nil
			})
			if err != nil {
				t.Fatalf("Stream: %v", err)
			}
			if reason != "end_turn" {
				t.Fatalf("StopReason = %q, want %q", reason, "end_turn")
			}
			if got := strings.Join(deltas, ""); !strings.Contains(got, tt.wantMarker) {
				t.Fatalf("permission marker = %q, want contains %q", got, tt.wantMarker)
			}
		})
	}
}

// TestBlackboxE2ESmoke performs a real turn with the installed blackbox binary.
// Run with: E2E_BLACKBOX=1 go test ./internal/agents/blackbox -run TestBlackboxE2ESmoke -v -timeout 120s
func TestBlackboxE2ESmoke(t *testing.T) {
	if os.Getenv("E2E_BLACKBOX") != "1" {
		t.Skip("set E2E_BLACKBOX=1 to run")
	}
	if err := blackbox.Preflight(); err != nil {
		t.Skipf("blackbox not available: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	client, err := blackbox.New(blackbox.Config{Dir: cwd})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	var builder strings.Builder
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	reason, err := client.Stream(ctx, "Reply with exactly the word PONG and nothing else.", func(delta string) error {
		fmt.Print(delta)
		builder.WriteString(delta)
		return nil
	})
	fmt.Println()

	if err != nil {
		t.Fatalf("Stream: %v", err)
	}
	if reason != "end_turn" {
		t.Fatalf("StopReason = %q, want %q", reason, "end_turn")
	}
	if builder.Len() == 0 {
		t.Fatal("no response text received")
	}
}
