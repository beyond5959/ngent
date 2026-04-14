package droid

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/agents/acpcli"
	"github.com/beyond5959/ngent/internal/agents/acpstdio"
	"github.com/beyond5959/ngent/internal/agents/agentutil"
)

var handlePermissionRequest = acpcli.StructuredPermissionRequestHandler(acpcli.DefaultPermissionTimeout)

// Config configures the Factory Droid ACP stdio provider.
type Config = agentutil.Config

// Client runs one Factory Droid ACP process per ACP operation.
type Client struct {
	*acpcli.Client
}

var _ agents.Streamer = (*Client)(nil)
var _ agents.ConfigOptionManager = (*Client)(nil)
var _ agents.SessionLister = (*Client)(nil)
var _ agents.SessionTranscriptLoader = (*Client)(nil)
var _ agents.SlashCommandsProvider = (*Client)(nil)

// New constructs a Factory Droid ACP client.
func New(cfg Config) (*Client, error) {
	base, err := acpcli.New(agents.AgentIDDroid, cfg, acpcli.Hooks{
		OpenConn:                openConn(cfg.Dir),
		SessionNewParams:        sessionNewParams(cfg.Dir),
		SessionLoadParams:       acpcli.SessionLoadParams(cfg.Dir),
		SessionListParams:       acpcli.SessionListParams(cfg.Dir),
		PromptParams:            promptParams,
		DiscoverModelsParams:    sessionNewParams(cfg.Dir),
		PrepareConfigSession:    prepareConfigSession,
		HandlePermissionRequest: handlePermissionRequest,
	})
	if err != nil {
		return nil, err
	}
	return &Client{Client: base}, nil
}

// Preflight checks that the droid binary is available in PATH.
func Preflight() error {
	return agentutil.PreflightBinary(agents.AgentIDDroid)
}

func openConn(dir string) func(context.Context, acpcli.OpenConnRequest) (*acpstdio.Conn, func(), json.RawMessage, error) {
	return func(
		ctx context.Context,
		req acpcli.OpenConnRequest,
	) (*acpstdio.Conn, func(), json.RawMessage, error) {
		modelID := strings.TrimSpace(req.ModelID)
		switch req.Purpose {
		case acpcli.OpenPurposeDiscoverModels, acpcli.OpenPurposeSessionList, acpcli.OpenPurposeTranscript:
			modelID = ""
		}

		conn, cleanup, initResult, err := acpcli.OpenProcess(ctx, acpcli.ProcessConfig{
			Command: agents.AgentIDDroid,
			Args:    commandArgs(modelID),
			Dir:     strings.TrimSpace(dir),
			Env:     os.Environ(),
			ConnOptions: acpstdio.ConnOptions{
				Prefix:           agents.AgentIDDroid,
				AllowStdoutNoise: true,
			},
			InitializeParams: initializeParams(),
		})
		if err != nil {
			return nil, nil, nil, acpcli.WrapOpenError(agents.AgentIDDroid, req.Purpose, err)
		}
		return conn, cleanup, initResult, nil
	}
}

func commandArgs(modelID string) []string {
	args := []string{"exec", "--output-format", "acp"}
	if modelID = strings.TrimSpace(modelID); modelID != "" {
		args = append(args, "--model", modelID)
	}
	return args
}

func initializeParams() map[string]any {
	return map[string]any{
		"protocolVersion": 1,
		"clientInfo": map[string]any{
			"name":    "ngent",
			"version": "0.1.0",
		},
		"clientCapabilities": map[string]any{
			"fs": map[string]any{
				"readTextFile":  false,
				"writeTextFile": false,
			},
		},
	}
}

func sessionNewParams(dir string) func(string) map[string]any {
	return func(string) map[string]any {
		return map[string]any{
			"cwd":        strings.TrimSpace(dir),
			"mcpServers": []any{},
		}
	}
}

func promptParams(sessionID string, prompt agents.Prompt, _ string) map[string]any {
	return acpcli.ACPPromptParams(sessionID, prompt)
}

func prepareConfigSession(
	modelID string,
	_ map[string]string,
	configID, value string,
) acpcli.ConfigSessionPlan {
	plan := acpcli.ConfigSessionPlan{
		SessionModelID: strings.TrimSpace(modelID),
	}
	if strings.EqualFold(strings.TrimSpace(configID), "model") && strings.TrimSpace(value) != "" {
		plan.SessionModelID = strings.TrimSpace(value)
		plan.SkipSetConfig = true
	}
	return plan
}

// Name returns the provider identifier.
func (c *Client) Name() string {
	if c == nil || c.Client == nil {
		return agents.AgentIDDroid
	}
	return c.Client.Name()
}
