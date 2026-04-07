package pi

import (
	"context"
	"fmt"
	"strings"

	"github.com/beyond5959/acp-adapter/pkg/piacp"
	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/agents/acpmodel"
)

// DiscoverModels queries ACP session/new and returns selectable model options.
func DiscoverModels(ctx context.Context, cfg Config) ([]agents.ModelOption, error) {
	cfg.Dir = strings.TrimSpace(cfg.Dir)
	if cfg.Dir == "" {
		return nil, fmt.Errorf("pi: discover models requires non-empty dir")
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

	runtime := piacp.NewEmbeddedRuntime(client.runtimeConfig)
	if err := runtime.Start(context.Background()); err != nil {
		return nil, fmt.Errorf("pi: start runtime for model discovery: %w", err)
	}
	defer runtime.Close()

	if _, err := client.clientRequest(requestCtx, runtime, methodInitialize, map[string]any{
		"client": map[string]any{
			"name": "ngent",
		},
	}); err != nil {
		return nil, fmt.Errorf("pi: initialize for model discovery: %w", err)
	}

	sessionResp, err := client.clientRequest(requestCtx, runtime, methodSessionNew, map[string]any{
		"cwd": client.Dir(),
	})
	if err != nil {
		return nil, fmt.Errorf("pi: session/new for model discovery: %w", err)
	}

	return acpmodel.ExtractModelOptions(sessionResp.Result), nil
}
