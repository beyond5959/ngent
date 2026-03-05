package codex

import (
	"context"
	"fmt"
	"strings"

	"github.com/beyond5959/acp-adapter/pkg/codexacp"
	"github.com/beyond5959/go-acp-server/internal/agents"
	"github.com/beyond5959/go-acp-server/internal/agents/acpmodel"
)

// DiscoverModels queries ACP session/new and returns selectable model options.
func DiscoverModels(ctx context.Context, cfg Config) ([]agents.ModelOption, error) {
	cfg.Dir = strings.TrimSpace(cfg.Dir)
	if cfg.Dir == "" {
		return nil, fmt.Errorf("codex: discover models requires non-empty dir")
	}
	cfg.ModelID = ""

	client, err := New(cfg)
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}

	requestCtx, cancel := context.WithTimeout(ctx, client.startTimeout)
	defer cancel()

	runtime := codexacp.NewEmbeddedRuntime(client.runtimeConfig)
	if err := runtime.Start(context.Background()); err != nil {
		return nil, fmt.Errorf("codex: start runtime for model discovery: %w", err)
	}
	defer runtime.Close()

	if _, err := client.clientRequest(requestCtx, runtime, methodInitialize, map[string]any{
		"client": map[string]any{
			"name": "go-acp-server",
		},
	}); err != nil {
		return nil, fmt.Errorf("codex: initialize for model discovery: %w", err)
	}

	sessionResp, err := client.clientRequest(requestCtx, runtime, methodSessionNew, map[string]any{
		"cwd": client.dir,
	})
	if err != nil {
		return nil, fmt.Errorf("codex: session/new for model discovery: %w", err)
	}

	return acpmodel.ExtractModelOptions(sessionResp.Result), nil
}
