package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	agentimpl "github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/runtime"
	"github.com/beyond5959/ngent/internal/storage"
)

func TestValidateListenAddr(t *testing.T) {
	tests := []struct {
		name        string
		listenAddr  string
		allowPublic bool
		wantErr     bool
		wantPort    int
	}{
		{
			name:        "loopback_allowed_when_public_disabled",
			listenAddr:  "127.0.0.1:8686",
			allowPublic: false,
			wantErr:     false,
			wantPort:    8686,
		},
		{
			name:        "localhost_allowed",
			listenAddr:  "localhost:8080",
			allowPublic: false,
			wantErr:     false,
			wantPort:    8080,
		},
		{
			name:        "public_ipv4_denied_when_public_disabled",
			listenAddr:  "0.0.0.0:8686",
			allowPublic: false,
			wantErr:     true,
		},
		{
			name:        "public_ipv6_denied_when_public_disabled",
			listenAddr:  "[::]:8686",
			allowPublic: false,
			wantErr:     true,
		},
		{
			name:        "public_ipv4_allowed_when_public_enabled",
			listenAddr:  "0.0.0.0:8686",
			allowPublic: true,
			wantErr:     false,
			wantPort:    8686,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotPort, err := validateListenAddr(tt.listenAddr, tt.allowPublic)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("validateListenAddr(%q, %v) error = nil, want non-nil", tt.listenAddr, tt.allowPublic)
				}
				return
			}

			if err != nil {
				t.Fatalf("validateListenAddr(%q, %v) unexpected error: %v", tt.listenAddr, tt.allowPublic, err)
			}
			if gotPort != tt.wantPort {
				t.Fatalf("port = %d, want %d", gotPort, tt.wantPort)
			}
		})
	}
}

func TestResolveAllowedRoots(t *testing.T) {
	roots, err := resolveAllowedRoots()
	if err != nil {
		t.Fatalf("resolveAllowedRoots() unexpected error: %v", err)
	}
	if got, want := len(roots), 1; got != want {
		t.Fatalf("len(roots) = %d, want %d", got, want)
	}
	if !filepath.IsAbs(roots[0]) {
		t.Fatalf("root %q is not absolute", roots[0])
	}
}

func TestResolveModelDiscoveryDir(t *testing.T) {
	root := t.TempDir()
	if got := resolveModelDiscoveryDir([]string{root}); got != root {
		t.Fatalf("resolveModelDiscoveryDir() = %q, want %q", got, root)
	}

	t.Run("fallback to cwd when roots missing", func(t *testing.T) {
		got := resolveModelDiscoveryDir([]string{filepath.Join(root, "missing")})
		if strings.TrimSpace(got) == "" {
			t.Fatalf("resolveModelDiscoveryDir() returned empty path")
		}
		if !filepath.IsAbs(got) {
			t.Fatalf("resolveModelDiscoveryDir() = %q, want absolute path", got)
		}
	})
}

func TestExtractConfigOverrides(t *testing.T) {
	got := extractConfigOverrides(`{
		"modelId":"gpt-5",
		"configOverrides":{
			"thought_level":"high",
			" empty ":" ",
			"non_string": 1
		}
	}`)
	if len(got) != 1 {
		t.Fatalf("len(configOverrides) = %d, want 1", len(got))
	}
	if got["thought_level"] != "high" {
		t.Fatalf("thought_level = %q, want %q", got["thought_level"], "high")
	}
}

func TestAgentConfigCatalogRefresherRefresh(t *testing.T) {
	store := newCatalogTestStore(t)
	defer func() {
		_ = store.Close()
	}()

	refresher := &agentConfigCatalogRefresher{
		store:    store,
		logger:   slog.New(slog.NewJSONHandler(io.Discard, nil)),
		agentIDs: []string{"codex"},
		fetchConfigOptions: func(ctx context.Context, agentID, modelID string) ([]agentimpl.ConfigOption, error) {
			_ = ctx
			_ = agentID
			switch modelID {
			case "":
				return []agentimpl.ConfigOption{
					modelCatalogOption("gpt-5", "gpt-5", "GPT-5", "gpt-5-mini", "GPT-5 Mini"),
					reasoningCatalogOption("medium", "low", "medium", "high"),
				}, nil
			case "gpt-5":
				return []agentimpl.ConfigOption{
					modelCatalogOption("gpt-5", "gpt-5", "GPT-5", "gpt-5-mini", "GPT-5 Mini"),
					reasoningCatalogOption("high", "medium", "high"),
				}, nil
			case "gpt-5-mini":
				return []agentimpl.ConfigOption{
					modelCatalogOption("gpt-5-mini", "gpt-5-mini", "GPT-5 Mini", "gpt-5", "GPT-5"),
					reasoningCatalogOption("low", "low", "medium"),
				}, nil
			default:
				return nil, context.DeadlineExceeded
			}
		},
		discoverModels: func(ctx context.Context, agentID string, defaultOptions []agentimpl.ConfigOption) ([]agentimpl.ModelOption, error) {
			_ = ctx
			_ = agentID
			return modelOptionsFromConfigOptions(defaultOptions), nil
		},
	}

	refresher.Refresh(context.Background())

	catalogs, err := store.ListAgentConfigCatalogsByAgent(context.Background(), "codex")
	if err != nil {
		t.Fatalf("ListAgentConfigCatalogsByAgent(): %v", err)
	}
	if got, want := len(catalogs), 3; got != want {
		t.Fatalf("len(catalogs) = %d, want %d", got, want)
	}

	defaultOptions := loadCatalogOptions(t, store, "codex", storage.DefaultAgentConfigCatalogModelID)
	if got := currentConfigValue(defaultOptions, "model"); got != "gpt-5" {
		t.Fatalf("default model currentValue = %q, want %q", got, "gpt-5")
	}
	miniOptions := loadCatalogOptions(t, store, "codex", "gpt-5-mini")
	if got := currentConfigValue(miniOptions, "reasoning"); got != "low" {
		t.Fatalf("mini reasoning currentValue = %q, want %q", got, "low")
	}
}

func TestAgentConfigCatalogRefresherPartialKeepsExistingRows(t *testing.T) {
	store := newCatalogTestStore(t)
	defer func() {
		_ = store.Close()
	}()

	if err := store.ReplaceAgentConfigCatalogs(context.Background(), "codex", []storage.UpsertAgentConfigCatalogParams{
		mustCatalogEntry(t, "codex", storage.DefaultAgentConfigCatalogModelID, []agentimpl.ConfigOption{
			modelCatalogOption("gpt-5", "gpt-5", "GPT-5", "gpt-5-mini", "GPT-5 Mini"),
			reasoningCatalogOption("medium", "low", "medium"),
		}),
		mustCatalogEntry(t, "codex", "gpt-5", []agentimpl.ConfigOption{
			modelCatalogOption("gpt-5", "gpt-5", "GPT-5", "gpt-5-mini", "GPT-5 Mini"),
			reasoningCatalogOption("medium", "medium", "high"),
		}),
		mustCatalogEntry(t, "codex", "gpt-5-mini", []agentimpl.ConfigOption{
			modelCatalogOption("gpt-5-mini", "gpt-5-mini", "GPT-5 Mini", "gpt-5", "GPT-5"),
			reasoningCatalogOption("low", "low", "medium"),
		}),
	}); err != nil {
		t.Fatalf("ReplaceAgentConfigCatalogs(seed): %v", err)
	}

	refresher := &agentConfigCatalogRefresher{
		store:    store,
		logger:   slog.New(slog.NewJSONHandler(io.Discard, nil)),
		agentIDs: []string{"codex"},
		fetchConfigOptions: func(ctx context.Context, agentID, modelID string) ([]agentimpl.ConfigOption, error) {
			_ = ctx
			_ = agentID
			switch modelID {
			case "":
				return []agentimpl.ConfigOption{
					modelCatalogOption("gpt-5", "gpt-5", "GPT-5", "gpt-5-mini", "GPT-5 Mini"),
					reasoningCatalogOption("high", "medium", "high"),
				}, nil
			case "gpt-5":
				return []agentimpl.ConfigOption{
					modelCatalogOption("gpt-5", "gpt-5", "GPT-5", "gpt-5-mini", "GPT-5 Mini"),
					reasoningCatalogOption("high", "medium", "high"),
				}, nil
			case "gpt-5-mini":
				return nil, context.DeadlineExceeded
			default:
				return nil, context.DeadlineExceeded
			}
		},
		discoverModels: func(ctx context.Context, agentID string, defaultOptions []agentimpl.ConfigOption) ([]agentimpl.ModelOption, error) {
			_ = ctx
			_ = agentID
			return modelOptionsFromConfigOptions(defaultOptions), nil
		},
	}

	refresher.Refresh(context.Background())

	catalogs, err := store.ListAgentConfigCatalogsByAgent(context.Background(), "codex")
	if err != nil {
		t.Fatalf("ListAgentConfigCatalogsByAgent(): %v", err)
	}
	if got, want := len(catalogs), 3; got != want {
		t.Fatalf("len(catalogs) after partial refresh = %d, want %d", got, want)
	}

	defaultOptions := loadCatalogOptions(t, store, "codex", storage.DefaultAgentConfigCatalogModelID)
	if got := currentConfigValue(defaultOptions, "reasoning"); got != "high" {
		t.Fatalf("default reasoning currentValue = %q, want %q", got, "high")
	}
	miniOptions := loadCatalogOptions(t, store, "codex", "gpt-5-mini")
	if got := currentConfigValue(miniOptions, "reasoning"); got != "low" {
		t.Fatalf("mini reasoning currentValue = %q, want %q", got, "low")
	}
}

func TestSupportedAgentsCodexStatus(t *testing.T) {
	agentsUnavailable := supportedAgents(false, false, false, false, false, false)
	if len(agentsUnavailable) == 0 {
		t.Fatalf("supportedAgents returned empty list")
	}
	if agentsUnavailable[0].ID != "codex" {
		t.Fatalf("agents[0].ID = %q, want %q", agentsUnavailable[0].ID, "codex")
	}
	if agentsUnavailable[0].Status != "unavailable" {
		t.Fatalf("codex unavailable status = %q, want %q", agentsUnavailable[0].Status, "unavailable")
	}
	if agentsUnavailable[1].ID != "claude" {
		t.Fatalf("agents[1].ID = %q, want %q", agentsUnavailable[1].ID, "claude")
	}
	if agentsUnavailable[1].Status != "unavailable" {
		t.Fatalf("claude unavailable status = %q, want %q", agentsUnavailable[1].Status, "unavailable")
	}
	if got, want := len(agentsUnavailable), 6; got != want {
		t.Fatalf("len(agentsUnavailable) = %d, want %d", got, want)
	}
	if agentsUnavailable[3].ID != "kimi" {
		t.Fatalf("agents[3].ID = %q, want %q", agentsUnavailable[3].ID, "kimi")
	}
	if agentsUnavailable[3].Status != "unavailable" {
		t.Fatalf("kimi unavailable status = %q, want %q", agentsUnavailable[3].Status, "unavailable")
	}
	if agentsUnavailable[4].ID != "qwen" {
		t.Fatalf("agents[4].ID = %q, want %q", agentsUnavailable[4].ID, "qwen")
	}
	if agentsUnavailable[4].Status != "unavailable" {
		t.Fatalf("qwen unavailable status = %q, want %q", agentsUnavailable[4].Status, "unavailable")
	}
	if agentsUnavailable[5].ID != "opencode" {
		t.Fatalf("agents[5].ID = %q, want %q", agentsUnavailable[5].ID, "opencode")
	}
	if agentsUnavailable[5].Status != "unavailable" {
		t.Fatalf("opencode unavailable status = %q, want %q", agentsUnavailable[5].Status, "unavailable")
	}

	agentsAvailable := supportedAgents(true, true, true, true, true, true)
	if agentsAvailable[0].Status != "available" {
		t.Fatalf("codex available status = %q, want %q", agentsAvailable[0].Status, "available")
	}
	if agentsAvailable[1].ID != "claude" {
		t.Fatalf("agents[1].ID = %q, want %q", agentsAvailable[1].ID, "claude")
	}
	if agentsAvailable[1].Status != "available" {
		t.Fatalf("claude available status = %q, want %q", agentsAvailable[1].Status, "available")
	}
	if got, want := len(agentsAvailable), 6; got != want {
		t.Fatalf("len(agentsAvailable) = %d, want %d", got, want)
	}
	if agentsAvailable[3].ID != "kimi" {
		t.Fatalf("agents[3].ID = %q, want %q", agentsAvailable[3].ID, "kimi")
	}
	if agentsAvailable[3].Status != "available" {
		t.Fatalf("kimi available status = %q, want %q", agentsAvailable[3].Status, "available")
	}
	if agentsAvailable[4].ID != "qwen" {
		t.Fatalf("agents[4].ID = %q, want %q", agentsAvailable[4].ID, "qwen")
	}
	if agentsAvailable[4].Status != "available" {
		t.Fatalf("qwen available status = %q, want %q", agentsAvailable[4].Status, "available")
	}
	if agentsAvailable[5].ID != "opencode" {
		t.Fatalf("agents[5].ID = %q, want %q", agentsAvailable[5].ID, "opencode")
	}
	if agentsAvailable[5].Status != "available" {
		t.Fatalf("opencode available status = %q, want %q", agentsAvailable[5].Status, "available")
	}
}

func newCatalogTestStore(t *testing.T) *storage.Store {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "catalog.db")
	store, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("storage.New(%q): %v", dbPath, err)
	}
	return store
}

func modelCatalogOption(current string, optionOneID string, optionOneName string, optionTwoID string, optionTwoName string) agentimpl.ConfigOption {
	return agentimpl.ConfigOption{
		ID:           "model",
		Category:     "model",
		Type:         "select",
		CurrentValue: current,
		Options: []agentimpl.ConfigOptionValue{
			{Value: optionOneID, Name: optionOneName},
			{Value: optionTwoID, Name: optionTwoName},
		},
	}
}

func reasoningCatalogOption(current string, values ...string) agentimpl.ConfigOption {
	options := make([]agentimpl.ConfigOptionValue, 0, len(values))
	for _, value := range values {
		options = append(options, agentimpl.ConfigOptionValue{Value: value, Name: value})
	}
	return agentimpl.ConfigOption{
		ID:           "reasoning",
		Category:     "reasoning",
		Type:         "select",
		CurrentValue: current,
		Options:      options,
	}
}

func mustCatalogEntry(
	t *testing.T,
	agentID string,
	modelID string,
	options []agentimpl.ConfigOption,
) storage.UpsertAgentConfigCatalogParams {
	t.Helper()

	entry, err := newAgentConfigCatalogEntry(agentID, modelID, options)
	if err != nil {
		t.Fatalf("newAgentConfigCatalogEntry(): %v", err)
	}
	return entry
}

func loadCatalogOptions(t *testing.T, store *storage.Store, agentID, modelID string) []agentimpl.ConfigOption {
	t.Helper()

	catalog, err := store.GetAgentConfigCatalog(context.Background(), agentID, modelID)
	if err != nil {
		t.Fatalf("GetAgentConfigCatalog(%q, %q): %v", agentID, modelID, err)
	}

	var options []agentimpl.ConfigOption
	if err := json.Unmarshal([]byte(catalog.ConfigOptionsJSON), &options); err != nil {
		t.Fatalf("json.Unmarshal(catalog): %v", err)
	}
	return options
}

func currentConfigValue(options []agentimpl.ConfigOption, configID string) string {
	for _, option := range options {
		if option.ID == configID {
			return option.CurrentValue
		}
	}
	return ""
}

func TestResolveDefaultDBPath(t *testing.T) {
	const home = "/tmp/test-home-db-default"
	t.Setenv("HOME", home)

	got, err := resolveDefaultDBPath()
	if err != nil {
		t.Fatalf("resolveDefaultDBPath() unexpected error: %v", err)
	}

	want := filepath.Join(home, ".go-agent-server", "agent-hub.db")
	if got != want {
		t.Fatalf("resolveDefaultDBPath() = %q, want %q", got, want)
	}
}

func TestEnsureDBPathParent(t *testing.T) {
	t.Run("create nested parent dir", func(t *testing.T) {
		tmp := t.TempDir()
		dbPath := filepath.Join(tmp, "nested", "dir", "agent-hub.db")
		if err := ensureDBPathParent(dbPath); err != nil {
			t.Fatalf("ensureDBPathParent(%q) unexpected error: %v", dbPath, err)
		}

		parent := filepath.Dir(dbPath)
		info, err := os.Stat(parent)
		if err != nil {
			t.Fatalf("os.Stat(%q): %v", parent, err)
		}
		if !info.IsDir() {
			t.Fatalf("parent %q is not a directory", parent)
		}
	})

	t.Run("reject empty path", func(t *testing.T) {
		if err := ensureDBPathParent("   "); err == nil {
			t.Fatalf("ensureDBPathParent should fail for empty path")
		}
	})
}

func TestGracefulShutdownForceCancelsTurns(t *testing.T) {
	controller := runtime.NewTurnController()
	cancelled := make(chan struct{}, 1)
	cancelFn := func() {
		select {
		case cancelled <- struct{}{}:
		default:
		}
	}

	if err := controller.Activate("th-1", "tu-1", cancelFn); err != nil {
		t.Fatalf("Activate() unexpected error: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	gracefulShutdown(context.Background(), logger, &http.Server{}, controller, 50*time.Millisecond)

	select {
	case <-cancelled:
	case <-time.After(2 * time.Second):
		t.Fatalf("turn cancel function was not called")
	}
}

func TestGetLANURLReturnsFalseForLoopback(t *testing.T) {
	url, ok := getLANURL("127.0.0.1:8686")
	if ok {
		t.Fatalf("getLANURL should return false for loopback")
	}
	if url != "" {
		t.Fatalf("expected empty URL for loopback, got %q", url)
	}
}

func TestPrintQRCodeDoesNothingForEmptyURL(t *testing.T) {
	var out bytes.Buffer
	printQRCode(&out, "")
	if got := out.String(); got != "" {
		t.Fatalf("printQRCode should write nothing for empty URL, got:\n%s", got)
	}
}
