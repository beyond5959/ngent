package claude

import (
	"context"
	"fmt"
	"strings"

	"github.com/beyond5959/acp-adapter/pkg/claudeacp"
	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/agents/acpmodel"
)

// DiscoverModels queries ACP session/new and returns selectable model options.
func DiscoverModels(ctx context.Context, cfg Config) ([]agents.ModelOption, error) {
	cfg.Dir = strings.TrimSpace(cfg.Dir)
	if cfg.Dir == "" {
		return nil, fmt.Errorf("claude: discover models requires non-empty dir")
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

	runtime := claudeacp.NewEmbeddedRuntime(client.runtimeConfig)
	if err := runtime.Start(context.Background()); err != nil {
		return nil, fmt.Errorf("claude: start runtime for model discovery: %w", err)
	}
	defer runtime.Close()

	if _, err := client.clientRequest(requestCtx, runtime, methodInitialize, map[string]any{
		"client": map[string]any{
			"name": "ngent",
		},
	}); err != nil {
		return nil, fmt.Errorf("claude: initialize for model discovery: %w", err)
	}

	sessionResp, err := client.clientRequest(requestCtx, runtime, methodSessionNew, map[string]any{
		"cwd": client.dir,
	})
	if err != nil {
		return nil, fmt.Errorf("claude: session/new for model discovery: %w", err)
	}

	return acpmodel.ExtractModelOptions(sessionResp.Result), nil
}
