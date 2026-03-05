package acpmodel

import (
	"encoding/json"
	"testing"

	"github.com/beyond5959/go-acp-server/internal/agents"
)

func TestExtractModelOptionsFromConfigOptions(t *testing.T) {
	raw := json.RawMessage(`{
		"sessionId":"ses_1",
		"configOptions":[
			{
				"id":"model",
				"type":"select",
				"currentValue":"gpt-5",
				"options":[
					{"value":"gpt-5","name":"GPT-5"},
					{"value":"gpt-5-mini","name":"GPT-5 Mini"}
				]
			}
		]
	}`)

	got := ExtractModelOptions(raw)
	if len(got) != 2 {
		t.Fatalf("len(models) = %d, want 2", len(got))
	}
	if got[0].ID != "gpt-5" || got[0].Name != "GPT-5" {
		t.Fatalf("models[0] = %+v, want id=gpt-5 name=GPT-5", got[0])
	}
	if got[1].ID != "gpt-5-mini" || got[1].Name != "GPT-5 Mini" {
		t.Fatalf("models[1] = %+v, want id=gpt-5-mini name=GPT-5 Mini", got[1])
	}
}

func TestExtractConfigOptionsFromConfigOptions(t *testing.T) {
	raw := json.RawMessage(`{
		"sessionId":"ses_1",
		"configOptions":[
			{
				"id":"model",
				"category":"model",
				"name":"Model",
				"description":"Model used for this session",
				"type":"select",
				"currentValue":"gpt-5.3-codex",
				"options":[
					{"value":"gpt-5.3-codex","name":"gpt-5.3-codex","description":"Latest frontier agentic coding model."},
					{"value":"gpt-5.2-codex","name":"gpt-5.2-codex","description":"Frontier agentic coding model."}
				]
			},
			{
				"id":"thought_level",
				"category":"reasoning",
				"name":"Thought level",
				"type":"select",
				"currentValue":"medium",
				"options":[
					{"value":"low","name":"Low"},
					{"value":"medium","name":"Medium"}
				]
			}
		]
	}`)

	got := ExtractConfigOptions(raw)
	if len(got) != 2 {
		t.Fatalf("len(configOptions) = %d, want %d", len(got), 2)
	}
	model, ok := FindModelConfigOption(got)
	if !ok {
		t.Fatalf("missing model config option")
	}
	if model.CurrentValue != "gpt-5.3-codex" {
		t.Fatalf("model currentValue = %q, want %q", model.CurrentValue, "gpt-5.3-codex")
	}
	if len(model.Options) != 2 {
		t.Fatalf("len(model options) = %d, want %d", len(model.Options), 2)
	}
	if model.Options[0].Description != "Latest frontier agentic coding model." {
		t.Fatalf("model options[0].description = %q", model.Options[0].Description)
	}
}

func TestExtractConfigOptionsFallbackFromModelsContainer(t *testing.T) {
	raw := json.RawMessage(`{
		"sessionId":"ses_2",
		"models":{
			"currentModelId":"gemini-2.5-pro",
			"availableModels":[
				{"modelId":"gemini-2.5-pro","name":"Gemini 2.5 Pro"},
				{"modelId":"gemini-2.5-flash","name":"Gemini 2.5 Flash"}
			]
		}
	}`)

	got := ExtractConfigOptions(raw)
	if len(got) != 1 {
		t.Fatalf("len(configOptions) = %d, want %d", len(got), 1)
	}
	model, ok := FindModelConfigOption(got)
	if !ok {
		t.Fatalf("missing model config option")
	}
	if model.CurrentValue != "gemini-2.5-pro" {
		t.Fatalf("model currentValue = %q, want %q", model.CurrentValue, "gemini-2.5-pro")
	}
	if len(model.Options) != 2 {
		t.Fatalf("len(model options) = %d, want %d", len(model.Options), 2)
	}
	if model.Options[1].Value != "gemini-2.5-flash" {
		t.Fatalf("model options[1].value = %q, want %q", model.Options[1].Value, "gemini-2.5-flash")
	}
}

func TestExtractModelOptionsFromModelsContainer(t *testing.T) {
	raw := json.RawMessage(`{
		"sessionId":"ses_2",
		"models":{
			"currentModelId":"gemini-2.5-pro",
			"availableModels":[
				"gemini-2.5-pro",
				{"modelId":"gemini-2.5-flash","name":"Gemini 2.5 Flash"}
			]
		}
	}`)

	got := ExtractModelOptions(raw)
	if len(got) != 2 {
		t.Fatalf("len(models) = %d, want 2", len(got))
	}
	if got[0].ID != "gemini-2.5-pro" {
		t.Fatalf("models[0].id = %q, want %q", got[0].ID, "gemini-2.5-pro")
	}
	if got[1].ID != "gemini-2.5-flash" || got[1].Name != "Gemini 2.5 Flash" {
		t.Fatalf("models[1] = %+v, want id=gemini-2.5-flash name=Gemini 2.5 Flash", got[1])
	}
}

func TestExtractModelOptionsInvalidPayload(t *testing.T) {
	got := ExtractModelOptions(json.RawMessage(`{"bad"`))
	if len(got) != 0 {
		t.Fatalf("len(models) = %d, want 0", len(got))
	}
}

func TestNormalizeModelOptions(t *testing.T) {
	input := []agents.ModelOption{
		{ID: "  gpt-5  ", Name: ""},
		{ID: "gpt-5", Name: "Ignored duplicate"},
		{ID: "gpt-5-mini", Name: "GPT-5 Mini"},
		{ID: " ", Name: "empty"},
	}

	got := NormalizeModelOptions(input)
	if len(got) != 2 {
		t.Fatalf("len(models) = %d, want 2", len(got))
	}
	if got[0].ID != "gpt-5" || got[0].Name != "gpt-5" {
		t.Fatalf("models[0] = %+v, want id=gpt-5 name=gpt-5", got[0])
	}
	if got[1].ID != "gpt-5-mini" || got[1].Name != "GPT-5 Mini" {
		t.Fatalf("models[1] = %+v, want id=gpt-5-mini name=GPT-5 Mini", got[1])
	}
}

func TestCurrentValueForConfig(t *testing.T) {
	options := []agents.ConfigOption{
		{ID: "model", CurrentValue: "gpt-5"},
		{ID: "thought_level", CurrentValue: "medium"},
	}

	if got := CurrentValueForConfig(options, "model"); got != "gpt-5" {
		t.Fatalf("CurrentValueForConfig(model) = %q, want %q", got, "gpt-5")
	}
	if got := CurrentValueForConfig(options, "unknown"); got != "" {
		t.Fatalf("CurrentValueForConfig(unknown) = %q, want empty", got)
	}
}
