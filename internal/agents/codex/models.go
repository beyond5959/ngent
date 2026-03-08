package codex

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/agents/acpmodel"
)

var (
	discoveryMu     sync.Mutex
	discoveryClient *Client
	discoveryKey    string
)

// CloseDiscoveryClient closes the shared model-discovery client, if any.
func CloseDiscoveryClient() error {
	discoveryMu.Lock()
	client := discoveryClient
	discoveryClient = nil
	discoveryKey = ""
	discoveryMu.Unlock()

	if client != nil {
		return client.Close()
	}
	return nil
}

// DiscoverModels queries ACP session/new and returns selectable model options.
func DiscoverModels(ctx context.Context, cfg Config) ([]agents.ModelOption, error) {
	cfg.Dir = strings.TrimSpace(cfg.Dir)
	if cfg.Dir == "" {
		return nil, fmt.Errorf("codex: discover models requires non-empty dir")
	}
	cfg.ModelID = ""

	if ctx == nil {
		ctx = context.Background()
	}

	client, err := getOrCreateDiscoveryClient(cfg)
	if err != nil {
		return nil, err
	}

	models, err := discoverModelsFromClient(ctx, client)
	if err == nil {
		return models, nil
	}

	// Retry once with a fresh shared client in case the cached runtime became unhealthy.
	resetDiscoveryClient(client)
	client, newErr := getOrCreateDiscoveryClient(cfg)
	if newErr != nil {
		return nil, newErr
	}
	models, err = discoverModelsFromClient(ctx, client)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func discoverModelsFromClient(ctx context.Context, client *Client) ([]agents.ModelOption, error) {
	options, err := client.ConfigOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("codex: session/new for model discovery: %w", err)
	}
	models := modelOptionsFromConfigOptions(options)
	if len(models) == 0 {
		return nil, errors.New("codex: model discovery returned empty options")
	}
	return models, nil
}

func modelOptionsFromConfigOptions(options []agents.ConfigOption) []agents.ModelOption {
	modelOption, ok := acpmodel.FindModelConfigOption(options)
	if !ok {
		return nil
	}

	out := make([]agents.ModelOption, 0, len(modelOption.Options)+1)
	for _, option := range modelOption.Options {
		id := strings.TrimSpace(option.Value)
		if id == "" {
			continue
		}
		name := strings.TrimSpace(option.Name)
		if name == "" {
			name = id
		}
		out = append(out, agents.ModelOption{ID: id, Name: name})
	}

	if current := strings.TrimSpace(modelOption.CurrentValue); current != "" {
		out = append(out, agents.ModelOption{ID: current, Name: current})
	}
	return acpmodel.NormalizeModelOptions(out)
}

func discoveryClientKey(cfg Config) string {
	parts := []string{
		strings.TrimSpace(cfg.Dir),
		strings.TrimSpace(cfg.Name),
		strings.TrimSpace(cfg.RuntimeConfig.AppServerCommand),
		strings.Join(cfg.RuntimeConfig.AppServerArgs, "\x1f"),
		strings.TrimSpace(cfg.RuntimeConfig.DefaultProfile),
		strings.TrimSpace(cfg.RuntimeConfig.InitialAuthMode),
	}
	return strings.Join(parts, "\x1e")
}

func getOrCreateDiscoveryClient(cfg Config) (*Client, error) {
	key := discoveryClientKey(cfg)

	discoveryMu.Lock()
	defer discoveryMu.Unlock()
	if discoveryClient != nil && discoveryKey == key {
		return discoveryClient, nil
	}

	if discoveryClient != nil {
		_ = discoveryClient.Close()
		discoveryClient = nil
		discoveryKey = ""
	}

	client, err := New(cfg)
	if err != nil {
		return nil, err
	}
	discoveryClient = client
	discoveryKey = key
	return discoveryClient, nil
}

func resetDiscoveryClient(target *Client) {
	discoveryMu.Lock()
	defer discoveryMu.Unlock()
	if discoveryClient != target {
		return
	}
	_ = discoveryClient.Close()
	discoveryClient = nil
	discoveryKey = ""
}
