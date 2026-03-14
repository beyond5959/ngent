package qwen

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
	return client.DiscoverModels(ctx)
}
