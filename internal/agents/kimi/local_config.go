package kimi

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/agents/acpmodel"
)

const (
	reasoningConfigID      = "reasoning"
	reasoningValueEnabled  = "enabled"
	reasoningValueDisabled = "disabled"
)

type localConfig struct {
	DefaultModelID  string
	DefaultThinking bool
	Models          []localModel
}

type localModel struct {
	ID               string
	Name             string
	SupportsThinking bool
}

func loadLocalConfig() (localConfig, error) {
	path, err := localConfigPath()
	if err != nil {
		return localConfig{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return localConfig{}, err
	}
	defer file.Close()

	var cfg localConfig
	var current *localModel

	commitCurrent := func() {
		if current == nil {
			return
		}
		current.ID = strings.TrimSpace(current.ID)
		if current.ID != "" {
			current.Name = strings.TrimSpace(current.Name)
			if current.Name == "" {
				current.Name = current.ID
			}
			cfg.Models = append(cfg.Models, *current)
		}
		current = nil
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(stripLineComment(scanner.Text()))
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			commitCurrent()
			modelID, ok := parseModelSection(line)
			if ok {
				current = &localModel{ID: modelID}
			}
			continue
		}

		key, value, ok := splitAssignment(line)
		if !ok {
			continue
		}

		switch {
		case current == nil && key == "default_model":
			cfg.DefaultModelID = trimQuoted(value)
		case current == nil && key == "default_thinking":
			cfg.DefaultThinking = strings.EqualFold(strings.TrimSpace(value), "true")
		case current != nil && key == "model":
			current.Name = trimQuoted(value)
		case current != nil && key == "capabilities":
			current.SupportsThinking = parseCapabilities(value)["thinking"]
		}
	}
	if err := scanner.Err(); err != nil {
		return localConfig{}, fmt.Errorf("scan kimi config: %w", err)
	}
	commitCurrent()

	if len(cfg.Models) == 0 {
		return localConfig{}, errors.New("kimi local config has no models")
	}

	sort.SliceStable(cfg.Models, func(i, j int) bool {
		if cfg.DefaultModelID != "" {
			if cfg.Models[i].ID == cfg.DefaultModelID {
				return true
			}
			if cfg.Models[j].ID == cfg.DefaultModelID {
				return false
			}
		}
		return cfg.Models[i].ID < cfg.Models[j].ID
	})

	if strings.TrimSpace(cfg.DefaultModelID) == "" {
		cfg.DefaultModelID = cfg.Models[0].ID
	}
	return cfg, nil
}

func localConfigPath() (string, error) {
	if path := strings.TrimSpace(os.Getenv("KIMI_CONFIG_FILE")); path != "" {
		return path, nil
	}
	if home := strings.TrimSpace(os.Getenv("KIMI_HOME")); home != "" {
		return filepath.Join(home, "config.toml"), nil
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(userHome, ".kimi", "config.toml"), nil
}

func (cfg localConfig) ConfigOptions(modelID string, overrides map[string]string) []agents.ConfigOption {
	if len(cfg.Models) == 0 {
		return nil
	}

	selectedModelID := strings.TrimSpace(modelID)
	if selectedModelID == "" {
		selectedModelID = cfg.DefaultModelID
	}

	selectedModel, hasSelectedModel := cfg.modelByID(selectedModelID)
	if !hasSelectedModel {
		selectedModel = localModel{ID: selectedModelID, Name: selectedModelID}
	}

	modelValues := make([]agents.ConfigOptionValue, 0, len(cfg.Models))
	for _, model := range cfg.Models {
		modelValues = append(modelValues, agents.ConfigOptionValue{
			Value: model.ID,
			Name:  model.Name,
		})
	}

	options := []agents.ConfigOption{{
		ID:           "model",
		Category:     "model",
		Name:         "Model",
		Description:  "Model sourced from local Kimi CLI config",
		Type:         "select",
		CurrentValue: selectedModel.ID,
		Options:      modelValues,
	}}

	if !selectedModel.SupportsThinking {
		return acpmodel.NormalizeConfigOptions(options)
	}

	reasoningValue := cfg.defaultReasoningValue()
	if override, ok := normalizeThinkingValue(overrides[reasoningConfigID]); ok {
		reasoningValue = override
	}

	options = append(options, agents.ConfigOption{
		ID:           reasoningConfigID,
		Category:     reasoningConfigID,
		Name:         "Reasoning",
		Description:  "Whether Kimi thinking mode is enabled",
		Type:         "select",
		CurrentValue: reasoningValue,
		Options: []agents.ConfigOptionValue{
			{Value: reasoningValueEnabled, Name: "Enabled"},
			{Value: reasoningValueDisabled, Name: "Disabled"},
		},
	})

	return acpmodel.NormalizeConfigOptions(options)
}

func (cfg localConfig) ModelOptions() []agents.ModelOption {
	if len(cfg.Models) == 0 {
		return nil
	}

	options := make([]agents.ModelOption, 0, len(cfg.Models))
	for _, model := range cfg.Models {
		options = append(options, agents.ModelOption{
			ID:   model.ID,
			Name: model.Name,
		})
	}
	return acpmodel.NormalizeModelOptions(options)
}

func (cfg localConfig) SupportsThinking(modelID string) bool {
	model, ok := cfg.modelByID(modelID)
	return ok && model.SupportsThinking
}

func (cfg localConfig) modelByID(modelID string) (localModel, bool) {
	modelID = strings.TrimSpace(modelID)
	if modelID == "" {
		modelID = cfg.DefaultModelID
	}
	for _, model := range cfg.Models {
		if model.ID == modelID {
			return model, true
		}
	}
	return localModel{}, false
}

func (cfg localConfig) defaultReasoningValue() string {
	if cfg.DefaultThinking {
		return reasoningValueEnabled
	}
	return reasoningValueDisabled
}

func normalizeThinkingValue(value string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case reasoningValueEnabled, "true", "on", "thinking":
		return reasoningValueEnabled, true
	case reasoningValueDisabled, "false", "off", "standard":
		return reasoningValueDisabled, true
	default:
		return "", false
	}
}

func stripLineComment(line string) string {
	var b strings.Builder
	inQuote := false
	for i := 0; i < len(line); i++ {
		ch := line[i]
		if ch == '"' && (i == 0 || line[i-1] != '\\') {
			inQuote = !inQuote
		}
		if ch == '#' && !inQuote {
			break
		}
		b.WriteByte(ch)
	}
	return b.String()
}

func parseModelSection(line string) (string, bool) {
	section := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
	if !strings.HasPrefix(section, "models.") {
		return "", false
	}
	section = strings.TrimSpace(strings.TrimPrefix(section, "models."))
	if section == "" {
		return "", false
	}
	return trimQuoted(section), true
}

func splitAssignment(line string) (string, string, bool) {
	idx := strings.Index(line, "=")
	if idx <= 0 {
		return "", "", false
	}
	key := strings.TrimSpace(line[:idx])
	value := strings.TrimSpace(line[idx+1:])
	if key == "" {
		return "", "", false
	}
	return key, value, true
}

func trimQuoted(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
			return strings.TrimSpace(value[1 : len(value)-1])
		}
	}
	return value
}

func parseCapabilities(value string) map[string]bool {
	trimmed := strings.TrimSpace(value)
	if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") {
		return nil
	}
	trimmed = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "["), "]"))
	if trimmed == "" {
		return nil
	}

	capabilities := make(map[string]bool)
	for _, part := range strings.Split(trimmed, ",") {
		capability := strings.ToLower(strings.TrimSpace(trimQuoted(part)))
		if capability == "" {
			continue
		}
		capabilities[capability] = true
	}
	if len(capabilities) == 0 {
		return nil
	}
	return capabilities
}
