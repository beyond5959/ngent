package qwen

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/agents/acpmodel"
	"github.com/beyond5959/ngent/internal/agents/acpstdio"
)

// DiscoverModels starts one ACP session/new handshake and returns model options.
func DiscoverModels(ctx context.Context, cfg Config) ([]agents.ModelOption, error) {
	client, err := New(cfg)
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}

	cmd := exec.Command("qwen", "--acp")
	cmd.Dir = client.dir
	cmd.Env = os.Environ()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("qwen: discover models open stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("qwen: discover models open stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("qwen: discover models open stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("qwen: discover models start process: %w", err)
	}

	errCh := make(chan error, 1)
	go func() { _, _ = io.Copy(io.Discard, stderr) }()
	go func() { errCh <- cmd.Wait() }()

	conn := acpstdio.NewConn(stdin, stdout, "qwen")
	defer conn.Close()
	defer acpstdio.TerminateProcess(cmd, errCh, 2*time.Second)

	if _, err := conn.Call(ctx, "initialize", map[string]any{
		"protocolVersion": 1,
		"clientCapabilities": map[string]any{
			"fs": map[string]any{
				"readTextFile":  false,
				"writeTextFile": false,
			},
		},
	}); err != nil {
		return nil, fmt.Errorf("qwen: discover models initialize: %w", err)
	}

	newResult, err := conn.Call(ctx, "session/new", map[string]any{
		"cwd":        client.dir,
		"mcpServers": []any{},
	})
	if err != nil {
		return nil, fmt.Errorf("qwen: discover models session/new: %w", err)
	}

	return acpmodel.ExtractModelOptions(newResult), nil
}
