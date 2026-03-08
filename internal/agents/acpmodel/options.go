package acpmodel

import (
	"encoding/json"
	"strings"

	"github.com/beyond5959/ngent/internal/agents"
)

// ExtractConfigOptions parses ACP session config options from session/new or
// session/set_config_option results.
func ExtractConfigOptions(raw json.RawMessage) []agents.ConfigOption {
	if len(raw) == 0 {
		return nil
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}

	sessionOptions := NormalizeConfigOptions(parseSessionConfigOptions(payload["configOptions"]))
	if len(sessionOptions) > 0 {
		return sessionOptions
	}

	// Backward compatibility for providers that still expose model catalogs as
	// `models.availableModels` instead of ACP `configOptions`.
	modelOptions, current := parseModelsContainer(payload["models"])
	modelOptions = append(modelOptions, parseAvailableModels(payload["availableModels"])...)
	if model := trimString(payload["model"]); model != "" {
		modelOptions = append(modelOptions, agents.ModelOption{ID: model, Name: model})
		if current == "" {
			current = model
		}
	}

	modelOptions = NormalizeModelOptions(modelOptions)
	if len(modelOptions) == 0 {
		return nil
	}
	if current == "" {
		current = modelOptions[0].ID
	}

	values := make([]agents.ConfigOptionValue, 0, len(modelOptions))
	for _, option := range modelOptions {
		values = append(values, agents.ConfigOptionValue{
			Value: option.ID,
			Name:  option.Name,
		})
	}
	return []agents.ConfigOption{{
		ID:           "model",
		Category:     "model",
		Name:         "Model",
		Type:         "select",
		CurrentValue: current,
		Options:      values,
	}}
}

// ExtractModelOptions parses model picker options from ACP payloads.
func ExtractModelOptions(raw json.RawMessage) []agents.ModelOption {
	if len(raw) == 0 {
		return nil
	}

	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}

	sessionOptions := NormalizeConfigOptions(parseSessionConfigOptions(payload["configOptions"]))
	if modelConfig, ok := FindModelConfigOption(sessionOptions); ok {
		return modelOptionsFromConfig(modelConfig)
	}

	// Backward compatibility fallback.
	models, current := parseModelsContainer(payload["models"])
	models = append(models, parseAvailableModels(payload["availableModels"])...)
	if model := trimString(payload["model"]); model != "" {
		models = append(models, agents.ModelOption{ID: model, Name: model})
		if current == "" {
			current = model
		}
	}
	if current != "" {
		models = append(models, agents.ModelOption{ID: current, Name: current})
	}
	return NormalizeModelOptions(models)
}

// NormalizeModelOptions trims and deduplicates model options by id.
func NormalizeModelOptions(options []agents.ModelOption) []agents.ModelOption {
	if len(options) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(options))
	out := make([]agents.ModelOption, 0, len(options))
	for _, option := range options {
		id := strings.TrimSpace(option.ID)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		name := strings.TrimSpace(option.Name)
		if name == "" {
			name = id
		}
		seen[id] = struct{}{}
		out = append(out, agents.ModelOption{ID: id, Name: name})
	}
	return out
}

// NormalizeConfigOptions trims/deduplicates config options and their values.
func NormalizeConfigOptions(options []agents.ConfigOption) []agents.ConfigOption {
	if len(options) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(options))
	out := make([]agents.ConfigOption, 0, len(options))
	for _, option := range options {
		id := strings.TrimSpace(option.ID)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}

		values := normalizeConfigOptionValues(option.Options)
		current := strings.TrimSpace(option.CurrentValue)
		if current == "" && len(values) > 0 {
			current = values[0].Value
		}
		if current != "" && !containsConfigOptionValue(values, current) {
			values = append([]agents.ConfigOptionValue{{
				Value: current,
				Name:  current,
			}}, values...)
		}

		name := strings.TrimSpace(option.Name)
		if name == "" {
			name = id
		}
		configType := strings.TrimSpace(option.Type)
		if configType == "" && len(values) > 0 {
			configType = "select"
		}

		seen[id] = struct{}{}
		out = append(out, agents.ConfigOption{
			ID:           id,
			Category:     strings.TrimSpace(option.Category),
			Name:         name,
			Description:  strings.TrimSpace(option.Description),
			Type:         configType,
			CurrentValue: current,
			Options:      values,
		})
	}
	return out
}

// CloneConfigOptions deep-copies config options for safe sharing.
func CloneConfigOptions(options []agents.ConfigOption) []agents.ConfigOption {
	if len(options) == 0 {
		return nil
	}
	out := make([]agents.ConfigOption, 0, len(options))
	for _, option := range options {
		copied := option
		if len(option.Options) > 0 {
			copied.Options = append([]agents.ConfigOptionValue(nil), option.Options...)
		}
		out = append(out, copied)
	}
	return out
}

// FindModelConfigOption returns the model config option by category or id.
func FindModelConfigOption(options []agents.ConfigOption) (agents.ConfigOption, bool) {
	for _, option := range options {
		if strings.EqualFold(strings.TrimSpace(option.Category), "model") ||
			strings.EqualFold(strings.TrimSpace(option.ID), "model") {
			return option, true
		}
	}
	return agents.ConfigOption{}, false
}

// CurrentValueForConfig returns the current value for one config id.
func CurrentValueForConfig(options []agents.ConfigOption, configID string) string {
	configID = strings.TrimSpace(configID)
	if configID == "" {
		return ""
	}
	for _, option := range options {
		if strings.EqualFold(strings.TrimSpace(option.ID), configID) {
			return strings.TrimSpace(option.CurrentValue)
		}
	}
	return ""
}

func modelOptionsFromConfig(option agents.ConfigOption) []agents.ModelOption {
	out := make([]agents.ModelOption, 0, len(option.Options)+1)
	for _, value := range option.Options {
		id := strings.TrimSpace(value.Value)
		if id == "" {
			continue
		}
		name := strings.TrimSpace(value.Name)
		if name == "" {
			name = id
		}
		out = append(out, agents.ModelOption{ID: id, Name: name})
	}
	if current := strings.TrimSpace(option.CurrentValue); current != "" {
		out = append(out, agents.ModelOption{ID: current, Name: current})
	}
	return NormalizeModelOptions(out)
}

func parseSessionConfigOptions(value any) []agents.ConfigOption {
	items, ok := value.([]any)
	if !ok {
		return nil
	}

	out := make([]agents.ConfigOption, 0, len(items))
	for _, item := range items {
		entry, ok := item.(map[string]any)
		if !ok {
			continue
		}

		id := trimString(entry["id"])
		if id == "" {
			continue
		}

		out = append(out, agents.ConfigOption{
			ID:           id,
			Category:     trimString(entry["category"]),
			Name:         firstNonEmptyString(entry["name"], id),
			Description:  trimString(entry["description"]),
			Type:         trimString(entry["type"]),
			CurrentValue: trimString(entry["currentValue"]),
			Options:      parseConfigOptionValues(entry["options"]),
		})
	}
	return out
}

func parseConfigOptionValues(value any) []agents.ConfigOptionValue {
	items, ok := value.([]any)
	if !ok {
		return nil
	}

	out := make([]agents.ConfigOptionValue, 0, len(items))
	for _, item := range items {
		switch entry := item.(type) {
		case string:
			id := strings.TrimSpace(entry)
			if id == "" {
				continue
			}
			out = append(out, agents.ConfigOptionValue{
				Value: id,
				Name:  id,
			})
		case map[string]any:
			id := firstNonEmptyString(entry["value"], entry["modelId"], entry["id"], entry["name"])
			if id == "" {
				continue
			}
			out = append(out, agents.ConfigOptionValue{
				Value:       id,
				Name:        firstNonEmptyString(entry["name"], entry["title"], entry["label"], id),
				Description: trimString(entry["description"]),
			})
		}
	}
	return out
}

func parseModelsContainer(value any) ([]agents.ModelOption, string) {
	models, ok := value.(map[string]any)
	if !ok {
		return nil, ""
	}

	out := parseAvailableModels(models["availableModels"])
	currentID := firstNonEmptyString(
		models["currentModelId"],
		models["currentModel"],
		models["model"],
	)
	return out, currentID
}

func parseAvailableModels(value any) []agents.ModelOption {
	items, ok := value.([]any)
	if !ok {
		return nil
	}

	out := make([]agents.ModelOption, 0, len(items))
	for _, item := range items {
		switch entry := item.(type) {
		case string:
			id := strings.TrimSpace(entry)
			if id == "" {
				continue
			}
			out = append(out, agents.ModelOption{ID: id, Name: id})
		case map[string]any:
			id := firstNonEmptyString(entry["value"], entry["modelId"], entry["id"], entry["name"])
			if id == "" {
				continue
			}
			name := firstNonEmptyString(entry["name"], entry["title"], entry["label"], id)
			out = append(out, agents.ModelOption{ID: id, Name: name})
		}
	}
	return out
}

func normalizeConfigOptionValues(values []agents.ConfigOptionValue) []agents.ConfigOptionValue {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(values))
	out := make([]agents.ConfigOptionValue, 0, len(values))
	for _, value := range values {
		id := strings.TrimSpace(value.Value)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		name := strings.TrimSpace(value.Name)
		if name == "" {
			name = id
		}
		seen[id] = struct{}{}
		out = append(out, agents.ConfigOptionValue{
			Value:       id,
			Name:        name,
			Description: strings.TrimSpace(value.Description),
		})
	}
	return out
}

func containsConfigOptionValue(values []agents.ConfigOptionValue, target string) bool {
	target = strings.TrimSpace(target)
	if target == "" {
		return false
	}
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value.Value), target) {
			return true
		}
	}
	return false
}

func trimString(value any) string {
	text, _ := value.(string)
	return strings.TrimSpace(text)
}

func firstNonEmptyString(values ...any) string {
	for _, value := range values {
		if text := trimString(value); text != "" {
			return text
		}
	}
	return ""
}
