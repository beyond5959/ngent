package kimi

import (
	"context"

	"github.com/beyond5959/ngent/internal/agents"
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
	return client.DiscoverModels(ctx)
}
