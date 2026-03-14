package acpcli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/beyond5959/ngent/internal/agents/acpstdio"
)

// ProcessConfig describes one provider process launch.
type ProcessConfig struct {
	Command          string
	Args             []string
	Dir              string
	Env              []string
	ConnOptions      acpstdio.ConnOptions
	InitializeParams map[string]any
}

// OpenProcess starts one ACP CLI process, performs initialize, and returns the connection.
func OpenProcess(
	ctx context.Context,
	cfg ProcessConfig,
) (*acpstdio.Conn, func(), json.RawMessage, error) {
	command := strings.TrimSpace(cfg.Command)
	if command == "" {
		return nil, nil, nil, errorsf("command is required")
	}

	cmd := exec.Command(command, cfg.Args...)
	cmd.Dir = strings.TrimSpace(cfg.Dir)
	if len(cfg.Env) == 0 {
		cmd.Env = os.Environ()
	} else {
		cmd.Env = append([]string(nil), cfg.Env...)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, errorsf("open stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, errorsf("open stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, errorsf("open stderr pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, nil, errorsf("start process: %w", err)
	}

	errCh := make(chan error, 1)
	go func() { _, _ = io.Copy(io.Discard, stderr) }()
	go func() { errCh <- cmd.Wait() }()

	conn := acpstdio.NewConnWithOptions(stdin, stdout, cfg.ConnOptions)
	cleanup := func() {
		conn.Close()
		acpstdio.TerminateProcess(cmd, errCh, 2*time.Second)
	}

	initParams := cfg.InitializeParams
	if initParams == nil {
		initParams = map[string]any{}
	}
	initResult, err := conn.Call(ctx, "initialize", initParams)
	if err != nil {
		cleanup()
		return nil, nil, nil, errorsf("initialize: %w", err)
	}
	return conn, cleanup, initResult, nil
}

// WrapOpenError adds provider/purpose context to one process startup error.
func WrapOpenError(provider string, purpose OpenPurpose, err error) error {
	if err == nil {
		return nil
	}
	provider = strings.TrimSpace(provider)
	scope := strings.TrimSpace(string(purpose))
	switch {
	case provider == "":
		return err
	case scope == "", purpose == OpenPurposeStream:
		return fmt.Errorf("%s: %w", provider, err)
	default:
		return fmt.Errorf("%s: %s: %w", provider, scope, err)
	}
}

func errorsf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
