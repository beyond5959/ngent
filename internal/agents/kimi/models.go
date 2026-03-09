package kimi

import (
	"context"
	"fmt"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/agents/acpmodel"
)

// DiscoverModels starts one ACP session/new handshake and returns model options.
func DiscoverModels(ctx context.Context, cfg Config) ([]agents.ModelOption, error) {
	client, err := New(cfg)
	if err != nil {
		return nil, err
	}
	if localCfg, err := loadLocalConfig(); err == nil {
		return localCfg.ModelOptions(), nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	conn, cleanup, _, err := client.openConn(ctx, client.CurrentModelID(), client.CurrentConfigOverrides())
	if err != nil {
		return nil, err
	}
	defer cleanup()

	newResult, err := conn.Call(ctx, "session/new", sessionNewParams(client.Dir(), client.CurrentModelID()))
	if err != nil {
		return nil, fmt.Errorf("kimi: discover models session/new: %w", err)
	}

	return acpmodel.ExtractModelOptions(newResult), nil
}
