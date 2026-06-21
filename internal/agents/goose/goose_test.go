package goose_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	goose "github.com/beyond5959/ngent/internal/agents/goose"
)

// TestPreflight verifies that Preflight returns nil when the goose binary exists.
func TestPreflight(t *testing.T) {
	if _, err := exec.LookPath("goose"); err != nil {
		t.Skip("goose not in PATH")
	}
	if err := goose.Preflight(); err != nil {
		t.Fatalf("Preflight() = %v, want nil", err)
	}
}

// TestStreamWithFakeProcess tests the Stream protocol using a fake goose binary.
func TestStreamWithFakeProcess(t *testing.T) {
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
            "agentInfo":{"name":"goose","title":"Goose","version":"0.1.0"},
            "authMethods":[],
            "modes":{"currentModeId":"default","availableModes":[]},
            "agentCapabilities":{"loadSession":True}
        }})
    elif method == "session/new":
        send({"jsonrpc":"2.0","id":rid,"result":{
            "sessionId":"ses_goose_test_123",
            "models":{"currentModelId":"test-model","availableModels":[]}
        }})
    elif method == "session/prompt":
        sid = params.get("sessionId","")
        send({"jsonrpc":"2.0","method":"session/update","params":{
            "sessionId":sid,
            "update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"Hello from Goose"}}
        }})
        send({"jsonrpc":"2.0","id":rid,"result":{"stopReason":"end_turn"}})
        sys.exit(0)
    elif method == "session/cancel":
        send({"jsonrpc":"2.0","id":rid,"result":{}})
        sys.exit(0)
`, python3)

	tmpDir := t.TempDir()
	fakeBin := tmpDir + "/goose"
	if err := os.WriteFile(fakeBin, []byte(fakeScript), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	origPath := os.Getenv("PATH")
	t.Setenv("PATH", tmpDir+":"+origPath)

	if err := goose.Preflight(); err != nil {
		t.Fatalf("Preflight with fake binary: %v", err)
	}

	c, err := goose.New(goose.Config{Dir: tmpDir})
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
		t.Errorf("StopReason = %q, want %q", reason, "end_turn")
	}
	if len(deltas) == 0 {
		t.Error("no deltas received")
	}
	if got := strings.Join(deltas, ""); !strings.Contains(got, "Hello from Goose") {
		t.Errorf("deltas = %q, want to contain %q", got, "Hello from Goose")
	}
}

// TestDiscoverModelsWithFakeProcess tests model discovery.
func TestDiscoverModelsWithFakeProcess(t *testing.T) {
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
    if method == "initialize":
        send({"jsonrpc":"2.0","id":rid,"result":{
            "protocolVersion":1,
            "agentInfo":{"name":"goose","title":"Goose","version":"0.1.0"},
            "authMethods":[],
            "modes":{"currentModeId":"default","availableModes":[]},
            "agentCapabilities":{"loadSession":True}
        }})
    elif method == "session/new":
        send({"jsonrpc":"2.0","id":rid,"result":{
            "sessionId":"ses_goose_models_123",
            "models":{
                "currentModelId":"gpt-4o",
                "availableModels":[
                    "gpt-4o",
                    {"modelId":"gpt-4o-mini","name":"GPT-4o Mini"},
                    {"modelId":"claude-sonnet-4-5","name":"Claude Sonnet 4.5"}
                ]
            }
        }})
`, python3)

	tmpDir := t.TempDir()
	fakeBin := tmpDir + "/goose"
	if err := os.WriteFile(fakeBin, []byte(fakeScript), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	origPath := os.Getenv("PATH")
	t.Setenv("PATH", tmpDir+":"+origPath)

	models, err := goose.DiscoverModels(context.Background(), goose.Config{Dir: tmpDir})
	if err != nil {
		t.Fatalf("DiscoverModels: %v", err)
	}
	if got, want := len(models), 3; got != want {
		t.Fatalf("len(models) = %d, want %d", got, want)
	}
	if models[0].ID != "gpt-4o" {
		t.Fatalf("models[0].id = %q, want %q", models[0].ID, "gpt-4o")
	}
	if models[1].ID != "gpt-4o-mini" {
		t.Fatalf("models[1].id = %q, want %q", models[1].ID, "gpt-4o-mini")
	}
	if models[2].ID != "claude-sonnet-4-5" {
		t.Fatalf("models[2].id = %q, want %q", models[2].ID, "claude-sonnet-4-5")
	}
}

// TestGooseE2ESmoke performs a real turn with the installed goose binary.
// Run with: E2E_GOOSE=1 go test ./internal/agents/goose/ -run E2E -v -timeout 90s
func TestGooseE2ESmoke(t *testing.T) {
	if os.Getenv("E2E_GOOSE") != "1" {
		t.Skip("set E2E_GOOSE=1 to run")
	}
	if err := goose.Preflight(); err != nil {
		t.Skipf("goose not available: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	c, err := goose.New(goose.Config{Dir: cwd})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	var builder strings.Builder
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	reason, err := c.Stream(ctx, "Reply with exactly the word PONG and nothing else.", func(delta string) error {
		fmt.Print(delta)
		builder.WriteString(delta)
		return nil
	})
	fmt.Println()

	if err != nil {
		t.Fatalf("Stream: %v", err)
	}
	t.Logf("StopReason: %s", reason)
	t.Logf("Response: %q", builder.String())

	if reason != "end_turn" {
		t.Errorf("StopReason = %q, want %q", reason, "end_turn")
	}
	if builder.Len() == 0 {
		t.Error("no response text received")
	}
}
