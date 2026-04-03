package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/agents/acpmodel"
	"github.com/beyond5959/ngent/internal/storage"
)

type threadConfigSelectionState interface {
	CurrentModelID() string
	CurrentConfigOverrides() map[string]string
}

func (s *Server) handleAgentModels(w http.ResponseWriter, r *http.Request, agentID string) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, r)
		return
	}

	if _, ok := s.allowedAgent[agentID]; !ok {
		writeError(w, http.StatusNotFound, codeNotFound, "agent not found", map[string]any{
			"agent": agentID,
		})
		return
	}

	models, found, err := s.loadStoredAgentModels(r.Context(), agentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to load stored agent models", map[string]any{
			"agent":  agentID,
			"reason": err.Error(),
		})
		return
	}
	if !found {
		models = []agents.ModelOption{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"agentId": agentID,
		"models":  models,
	})
}

func (s *Server) handleThreadConfigOptions(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeMethodNotAllowed(w, r)
		return
	}

	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		options, found, err := s.loadStoredThreadConfigOptions(r.Context(), thread)
		if err != nil {
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to load stored thread config options", map[string]any{
				"threadId": thread.ThreadID,
				"reason":   err.Error(),
			})
			return
		}
		if !found {
			options = []agents.ConfigOption{}
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"threadId":      thread.ThreadID,
			"configOptions": options,
		})
	case http.MethodPost:
		if s.turns.IsThreadActive(thread.ThreadID) {
			writeError(w, http.StatusConflict, codeConflict, "thread has an active turn", map[string]any{"threadId": thread.ThreadID})
			return
		}

		var req struct {
			ConfigID string `json:"configId"`
			Value    string `json:"value"`
		}
		if err := decodeJSONBody(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "invalid JSON body", map[string]any{"reason": err.Error()})
			return
		}
		req.ConfigID = strings.TrimSpace(req.ConfigID)
		req.Value = strings.TrimSpace(req.Value)
		if req.ConfigID == "" {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "configId is required", map[string]any{"field": "configId"})
			return
		}
		if req.Value == "" {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "value is required", map[string]any{"field": "value"})
			return
		}

		currentOptions, err := s.loadThreadConfigOptionsForUpdate(r.Context(), thread)
		if err != nil {
			if errors.Is(err, errThreadConfigOptionsUnavailable) {
				writeError(w, http.StatusConflict, codeConflict, "thread config options are not available yet", map[string]any{
					"threadId": thread.ThreadID,
				})
				return
			}
			writeError(w, http.StatusServiceUnavailable, codeUpstreamUnavailable, "failed to load thread config options", map[string]any{
				"threadId": thread.ThreadID,
				"reason":   err.Error(),
			})
			return
		}
		if err := validateThreadConfigSelection(currentOptions, req.ConfigID, req.Value); err != nil {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, err.Error(), map[string]any{
				"threadId": thread.ThreadID,
				"configId": req.ConfigID,
				"value":    req.Value,
			})
			return
		}

		options, err := s.updatedThreadConfigOptions(r.Context(), thread, currentOptions, req.ConfigID, req.Value)
		if err != nil {
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to update thread config options", map[string]any{
				"threadId": thread.ThreadID,
				"reason":   err.Error(),
			})
			return
		}

		currentModel := acpmodel.CurrentValueForConfig(options, "model")
		agentOptionsJSON, err := withThreadConfigState(thread.AgentOptionsJSON, currentModel, options)
		if err != nil {
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to normalize thread agent options", map[string]any{
				"threadId": thread.ThreadID,
				"reason":   err.Error(),
			})
			return
		}
		if err := s.store.UpdateThreadAgentOptions(r.Context(), thread.ThreadID, agentOptionsJSON); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				writeThreadNotFound(w)
				return
			}
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to update thread", map[string]any{"reason": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"threadId":      thread.ThreadID,
			"configOptions": options,
		})
	}
}

func (s *Server) handleThreadSlashCommands(w http.ResponseWriter, r *http.Request, clientID, threadID string) {
	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	thread, ok := s.loadAccessibleThreadOrWriteNotFound(w, r.Context(), threadID)
	if !ok {
		return
	}

	commands, found, err := s.loadStoredAgentSlashCommands(r.Context(), thread.AgentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to load slash commands", map[string]any{
			"threadId": thread.ThreadID,
			"agent":    thread.AgentID,
			"reason":   err.Error(),
		})
		return
	}
	if !found {
		s.persistThreadSlashCommandsBestEffort(r.Context(), thread, nil)
		commands, found, err = s.loadStoredAgentSlashCommands(r.Context(), thread.AgentID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to load slash commands", map[string]any{
				"threadId": thread.ThreadID,
				"agent":    thread.AgentID,
				"reason":   err.Error(),
			})
			return
		}
		if !found {
			commands = []agents.SlashCommand{}
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"threadId": thread.ThreadID,
		"agentId":  thread.AgentID,
		"commands": commands,
	})
}

func (s *Server) loadStoredAgentModels(ctx context.Context, agentID string) ([]agents.ModelOption, bool, error) {
	catalogs, err := s.store.ListAgentConfigCatalogsByAgent(ctx, agentID)
	if err != nil {
		return nil, false, err
	}
	if len(catalogs) == 0 {
		return nil, false, nil
	}

	models := make([]agents.ModelOption, 0)
	for _, catalog := range catalogs {
		options, err := decodeStoredConfigOptions(catalog.ConfigOptionsJSON)
		if err != nil {
			s.logger.Warn("config_catalog.decode_failed",
				"agent", agentID,
				"modelId", catalog.ModelID,
				"reason", err.Error(),
			)
			continue
		}
		modelOption, ok := acpmodel.FindModelConfigOption(options)
		if !ok {
			continue
		}
		for _, value := range modelOption.Options {
			modelID := strings.TrimSpace(value.Value)
			if modelID == "" {
				continue
			}
			name := strings.TrimSpace(value.Name)
			if name == "" {
				name = modelID
			}
			models = append(models, agents.ModelOption{ID: modelID, Name: name})
		}
	}

	models = acpmodel.NormalizeModelOptions(models)
	if len(models) == 0 {
		return nil, false, nil
	}
	return models, true, nil
}

func (s *Server) loadStoredThreadConfigOptions(ctx context.Context, thread storage.Thread) ([]agents.ConfigOption, bool, error) {
	modelID, overrides := threadConfigSelections(thread.AgentOptionsJSON)
	if modelID != "" {
		catalog, err := s.store.GetAgentConfigCatalog(ctx, thread.AgentID, modelID)
		if errors.Is(err, storage.ErrNotFound) {
			return nil, false, nil
		}
		if err != nil {
			return nil, false, err
		}

		options, err := decodeStoredConfigOptions(catalog.ConfigOptionsJSON)
		if err != nil {
			return nil, false, err
		}
		return applyThreadConfigSelections(options, modelID, overrides), true, nil
	}

	sessionID := threadSessionID(thread.AgentOptionsJSON)
	if sessionID == "" {
		return nil, false, nil
	}

	return s.loadStoredSessionConfigOptions(ctx, thread, sessionID)
}

func (s *Server) loadStoredSessionConfigOptions(
	ctx context.Context,
	thread storage.Thread,
	sessionID string,
) ([]agents.ConfigOption, bool, error) {
	cache, err := s.store.GetSessionConfigCache(ctx, thread.AgentID, thread.CWD, sessionID)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	options, err := decodeStoredConfigOptions(cache.ConfigOptionsJSON)
	if err != nil {
		return nil, false, err
	}
	return options, true, nil
}

func (s *Server) loadStoredAgentConfigOptions(
	ctx context.Context,
	agentID, modelID string,
) ([]agents.ConfigOption, bool, error) {
	lookupModelID := strings.TrimSpace(modelID)
	if lookupModelID == "" {
		lookupModelID = storage.DefaultAgentConfigCatalogModelID
	}

	catalog, err := s.store.GetAgentConfigCatalog(ctx, agentID, lookupModelID)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	options, err := decodeStoredConfigOptions(catalog.ConfigOptionsJSON)
	if err != nil {
		return nil, false, err
	}
	return options, true, nil
}

func (s *Server) loadThreadConfigOptionsForUpdate(
	ctx context.Context,
	thread storage.Thread,
) ([]agents.ConfigOption, error) {
	options, found, err := s.loadStoredThreadConfigOptions(ctx, thread)
	if err != nil {
		return nil, err
	}
	if found {
		return options, nil
	}
	return nil, errThreadConfigOptionsUnavailable
}

func (s *Server) loadStoredAgentSlashCommands(ctx context.Context, agentID string) ([]agents.SlashCommand, bool, error) {
	stored, err := s.store.GetAgentSlashCommands(ctx, agentID)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	commands, err := decodeStoredSlashCommands(stored.CommandsJSON)
	if err != nil {
		return nil, false, err
	}
	return commands, true, nil
}

func (s *Server) persistThreadSlashCommandsBestEffort(ctx context.Context, thread storage.Thread, provider any) {
	if _, found, err := s.loadStoredAgentSlashCommands(ctx, thread.AgentID); err == nil && found {
		return
	}

	if provider == nil {
		resolved, err := s.resolveTurnAgent(thread)
		if err != nil {
			s.logger.Warn("thread.slash_commands_resolve_failed",
				"threadId", thread.ThreadID,
				"agent", thread.AgentID,
				"reason", err.Error(),
			)
			return
		}
		provider = resolved
	}

	reader, ok := provider.(agents.SlashCommandsProvider)
	if !ok {
		return
	}

	commands, known, err := reader.SlashCommands(ctx)
	if err != nil {
		s.logger.Warn("thread.slash_commands_probe_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"reason", err.Error(),
		)
		return
	}
	if !known {
		return
	}
	if err := s.persistAgentSlashCommands(ctx, thread.AgentID, commands); err != nil {
		s.logger.Warn("thread.slash_commands_persist_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"reason", err.Error(),
		)
	}
}

func (s *Server) persistAgentConfigCatalog(
	ctx context.Context,
	agentID string,
	agentOptionsJSON string,
	options []agents.ConfigOption,
) error {
	modelID, _ := threadConfigSelections(agentOptionsJSON)
	if currentModel := strings.TrimSpace(acpmodel.CurrentValueForConfig(options, "model")); currentModel != "" {
		modelID = currentModel
	}
	if strings.TrimSpace(modelID) == "" {
		modelID = storage.DefaultAgentConfigCatalogModelID
	}

	configOptionsJSON, err := encodeStoredConfigOptions(options)
	if err != nil {
		return err
	}

	return s.store.UpsertAgentConfigCatalog(ctx, storage.UpsertAgentConfigCatalogParams{
		AgentID:           agentID,
		ModelID:           modelID,
		ConfigOptionsJSON: configOptionsJSON,
	})
}

func (s *Server) updatedThreadConfigOptions(
	ctx context.Context,
	thread storage.Thread,
	currentOptions []agents.ConfigOption,
	configID, value string,
) ([]agents.ConfigOption, error) {
	modelID, overrides := threadConfigSelections(thread.AgentOptionsJSON)
	nextModelID := modelID
	nextOverrides := cloneThreadConfigOverrides(overrides)

	if strings.EqualFold(strings.TrimSpace(configID), "model") {
		nextModelID = strings.TrimSpace(value)
	} else {
		if nextOverrides == nil {
			nextOverrides = make(map[string]string)
		}
		nextOverrides[strings.TrimSpace(configID)] = strings.TrimSpace(value)
	}

	baseOptions := currentOptions
	if strings.EqualFold(strings.TrimSpace(configID), "model") {
		storedOptions, found, err := s.loadStoredAgentConfigOptions(ctx, thread.AgentID, nextModelID)
		if err != nil {
			return nil, err
		}
		if found {
			baseOptions = storedOptions
		} else {
			baseOptions = modelOnlyThreadConfigOptions(currentOptions)
			nextOverrides = nil
		}
	}

	return applyThreadConfigSelections(baseOptions, nextModelID, nextOverrides), nil
}

func (s *Server) syncThreadConfigSelections(
	ctx context.Context,
	thread storage.Thread,
	provider agents.Streamer,
) error {
	manager, ok := provider.(agents.ConfigOptionManager)
	if !ok {
		return nil
	}
	state, ok := provider.(threadConfigSelectionState)
	if !ok {
		return nil
	}

	desiredModelID, desiredOverrides := threadConfigSelections(thread.AgentOptionsJSON)
	currentModelID := strings.TrimSpace(state.CurrentModelID())
	if desiredModelID != "" && desiredModelID != currentModelID {
		options, err := manager.SetConfigOption(ctx, "model", desiredModelID)
		if err != nil {
			return fmt.Errorf("apply model config before turn: %w", err)
		}
		s.persistAgentConfigCatalogBestEffort(ctx, thread, options)
	}

	currentOverrides := state.CurrentConfigOverrides()
	if desiredModelID != "" && desiredModelID != currentModelID {
		currentOverrides = state.CurrentConfigOverrides()
	}
	configIDs := make([]string, 0, len(desiredOverrides))
	for configID := range desiredOverrides {
		configIDs = append(configIDs, configID)
	}
	sort.Strings(configIDs)

	for _, configID := range configIDs {
		value := strings.TrimSpace(desiredOverrides[configID])
		if value == "" || value == strings.TrimSpace(currentOverrides[configID]) {
			continue
		}
		options, err := manager.SetConfigOption(ctx, configID, value)
		if err != nil {
			return fmt.Errorf("apply %s config before turn: %w", configID, err)
		}
		s.persistAgentConfigCatalogBestEffort(ctx, thread, options)
	}

	return nil
}

func (s *Server) persistAgentConfigCatalogBestEffort(
	ctx context.Context,
	thread storage.Thread,
	options []agents.ConfigOption,
) {
	normalized := acpmodel.NormalizeConfigOptions(options)
	if len(normalized) == 0 {
		return
	}
	if err := s.persistAgentConfigCatalog(ctx, thread.AgentID, thread.AgentOptionsJSON, normalized); err != nil {
		s.logger.Warn("config_catalog.persist_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"reason", err.Error(),
		)
	}
}

func (s *Server) persistThreadConfigSnapshotBestEffort(
	ctx context.Context,
	thread *storage.Thread,
	options []agents.ConfigOption,
) {
	if thread == nil {
		return
	}

	normalized := acpmodel.NormalizeConfigOptions(options)
	if len(normalized) == 0 {
		return
	}

	currentModel := strings.TrimSpace(acpmodel.CurrentValueForConfig(normalized, "model"))
	nextAgentOptionsJSON, err := withThreadConfigState(thread.AgentOptionsJSON, currentModel, normalized)
	if err != nil {
		s.logger.Warn("thread.config_snapshot_encode_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"reason", err.Error(),
		)
		return
	}
	if nextAgentOptionsJSON != thread.AgentOptionsJSON {
		if err := s.store.UpdateThreadAgentOptions(ctx, thread.ThreadID, nextAgentOptionsJSON); err != nil {
			s.logger.Warn("thread.config_snapshot_persist_failed",
				"threadId", thread.ThreadID,
				"agent", thread.AgentID,
				"reason", err.Error(),
			)
			return
		}
		thread.AgentOptionsJSON = nextAgentOptionsJSON
	}
	if err := s.persistAgentConfigCatalog(ctx, thread.AgentID, thread.AgentOptionsJSON, normalized); err != nil {
		s.logger.Warn("config_catalog.persist_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"reason", err.Error(),
		)
	}
	s.persistSessionConfigSnapshotBestEffort(ctx, *thread, normalized)
}

func (s *Server) persistThreadSessionConfigSnapshotBestEffort(ctx context.Context, thread storage.Thread) {
	options, found, err := s.loadStoredThreadConfigOptions(ctx, thread)
	if err != nil {
		s.logger.Warn("thread.session_config_restore_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"sessionId", threadSessionID(thread.AgentOptionsJSON),
			"reason", err.Error(),
		)
		return
	}
	if !found {
		return
	}
	s.persistSessionConfigSnapshotBestEffort(ctx, thread, options)
}

func (s *Server) persistSessionConfigSnapshotBestEffort(
	ctx context.Context,
	thread storage.Thread,
	options []agents.ConfigOption,
) {
	s.persistSessionConfigSnapshotForSessionIDBestEffort(ctx, thread, threadSessionID(thread.AgentOptionsJSON), options)
}

func (s *Server) persistSessionConfigSnapshotForSessionIDBestEffort(
	ctx context.Context,
	thread storage.Thread,
	sessionID string,
	options []agents.ConfigOption,
) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return
	}

	configOptionsJSON, err := encodeStoredConfigOptions(options)
	if err != nil {
		s.logger.Warn("thread.session_config_encode_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"sessionId", sessionID,
			"reason", err.Error(),
		)
		return
	}

	if err := s.store.UpsertSessionConfigCache(ctx, storage.UpsertSessionConfigCacheParams{
		AgentID:           thread.AgentID,
		CWD:               thread.CWD,
		SessionID:         sessionID,
		ConfigOptionsJSON: configOptionsJSON,
	}); err != nil {
		s.logger.Warn("thread.session_config_persist_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"sessionId", sessionID,
			"reason", err.Error(),
		)
	}
}

func (s *Server) persistSessionLoadConfigSnapshotBestEffort(
	ctx context.Context,
	thread *storage.Thread,
	sessionID string,
	options []agents.ConfigOption,
) {
	if thread == nil {
		return
	}

	normalized := acpmodel.NormalizeConfigOptions(options)
	if len(normalized) == 0 {
		return
	}

	if threadSessionID(thread.AgentOptionsJSON) == strings.TrimSpace(sessionID) {
		s.persistThreadConfigSnapshotBestEffort(ctx, thread, normalized)
		return
	}

	if err := s.persistAgentConfigCatalog(ctx, thread.AgentID, thread.AgentOptionsJSON, normalized); err != nil {
		s.logger.Warn("config_catalog.persist_failed",
			"threadId", thread.ThreadID,
			"agent", thread.AgentID,
			"sessionId", sessionID,
			"reason", err.Error(),
		)
	}
	s.persistSessionConfigSnapshotForSessionIDBestEffort(ctx, *thread, sessionID, normalized)
}

func (s *Server) persistAgentSlashCommands(
	ctx context.Context,
	agentID string,
	commands []agents.SlashCommand,
) error {
	commandsJSON, err := encodeStoredSlashCommands(commands)
	if err != nil {
		return err
	}

	return s.store.UpsertAgentSlashCommands(ctx, storage.UpsertAgentSlashCommandsParams{
		AgentID:      agentID,
		CommandsJSON: commandsJSON,
	})
}

func normalizeAgentOptions(raw json.RawMessage) (string, error) {
	if len(strings.TrimSpace(string(raw))) == 0 {
		return "{}", nil
	}

	var objectValue map[string]any
	if err := json.Unmarshal(raw, &objectValue); err != nil {
		return "", err
	}

	normalized, err := json.Marshal(objectValue)
	if err != nil {
		return "", err
	}
	return string(normalized), nil
}

func withThreadConfigState(agentOptionsJSON, modelID string, options []agents.ConfigOption) (string, error) {
	modelID = strings.TrimSpace(modelID)

	objectValue := map[string]any{}
	trimmed := strings.TrimSpace(agentOptionsJSON)
	if trimmed != "" {
		if err := json.Unmarshal([]byte(trimmed), &objectValue); err != nil {
			return "", err
		}
	}

	if modelID == "" {
		delete(objectValue, "modelId")
	} else {
		objectValue["modelId"] = modelID
	}

	configOverrides := configOverridesFromOptions(options)
	if len(configOverrides) == 0 {
		delete(objectValue, "configOverrides")
	} else {
		objectValue["configOverrides"] = configOverrides
	}

	normalized, err := json.Marshal(objectValue)
	if err != nil {
		return "", err
	}
	return string(normalized), nil
}

func withoutThreadConfigState(agentOptionsJSON string) (string, error) {
	objectValue := map[string]any{}
	trimmed := strings.TrimSpace(agentOptionsJSON)
	if trimmed != "" {
		if err := json.Unmarshal([]byte(trimmed), &objectValue); err != nil {
			return "", err
		}
	}

	delete(objectValue, "modelId")
	delete(objectValue, "configOverrides")

	normalized, err := json.Marshal(objectValue)
	if err != nil {
		return "", err
	}
	return string(normalized), nil
}

func configOverridesFromOptions(options []agents.ConfigOption) map[string]string {
	overrides := make(map[string]string, len(options))
	for _, option := range options {
		configID := strings.TrimSpace(option.ID)
		if configID == "" || strings.EqualFold(configID, "model") {
			continue
		}
		value := strings.TrimSpace(option.CurrentValue)
		if value == "" {
			continue
		}
		overrides[configID] = value
	}
	if len(overrides) == 0 {
		return nil
	}
	return overrides
}

func cloneThreadConfigOverrides(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(input))
	for configID, value := range input {
		cloned[configID] = value
	}
	return cloned
}

func validateThreadConfigSelection(options []agents.ConfigOption, configID, value string) error {
	configID = strings.TrimSpace(configID)
	value = strings.TrimSpace(value)
	if configID == "" {
		return errors.New("configId is required")
	}
	if value == "" {
		return errors.New("value is required")
	}

	var option *agents.ConfigOption
	for i := range options {
		candidateID := strings.TrimSpace(options[i].ID)
		category := strings.TrimSpace(options[i].Category)
		if strings.EqualFold(candidateID, configID) {
			option = &options[i]
			break
		}
		if strings.EqualFold(configID, "model") && strings.EqualFold(category, "model") {
			option = &options[i]
			break
		}
	}
	if option == nil {
		return fmt.Errorf("config option %q is not available", configID)
	}
	if len(option.Options) == 0 {
		return nil
	}
	for _, candidate := range option.Options {
		if strings.EqualFold(strings.TrimSpace(candidate.Value), value) {
			return nil
		}
	}
	return fmt.Errorf("value %q is not available for config option %q", value, configID)
}

func decodeStoredConfigOptions(raw string) ([]agents.ConfigOption, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	var options []agents.ConfigOption
	if err := json.Unmarshal([]byte(raw), &options); err != nil {
		return nil, fmt.Errorf("decode stored config options: %w", err)
	}
	return acpmodel.NormalizeConfigOptions(options), nil
}

func encodeStoredConfigOptions(options []agents.ConfigOption) (string, error) {
	normalized := acpmodel.NormalizeConfigOptions(options)
	encoded, err := json.Marshal(normalized)
	if err != nil {
		return "", fmt.Errorf("encode stored config options: %w", err)
	}
	return string(encoded), nil
}

func decodeStoredSlashCommands(raw string) ([]agents.SlashCommand, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	var commands []agents.SlashCommand
	if err := json.Unmarshal([]byte(raw), &commands); err != nil {
		return nil, fmt.Errorf("decode stored slash commands: %w", err)
	}
	return agents.CloneSlashCommands(commands), nil
}

func encodeStoredSlashCommands(commands []agents.SlashCommand) (string, error) {
	normalized := agents.CloneSlashCommands(commands)
	encoded, err := json.Marshal(normalized)
	if err != nil {
		return "", fmt.Errorf("encode stored slash commands: %w", err)
	}
	return string(encoded), nil
}

func threadConfigSelections(agentOptionsJSON string) (string, map[string]string) {
	var raw struct {
		ModelID         string         `json:"modelId"`
		ConfigOverrides map[string]any `json:"configOverrides"`
	}
	if strings.TrimSpace(agentOptionsJSON) == "" {
		return "", nil
	}
	if err := json.Unmarshal([]byte(agentOptionsJSON), &raw); err != nil {
		return "", nil
	}

	overrides := make(map[string]string, len(raw.ConfigOverrides))
	for rawID, rawValue := range raw.ConfigOverrides {
		configID := strings.TrimSpace(rawID)
		if configID == "" {
			continue
		}
		value, ok := rawValue.(string)
		if !ok {
			continue
		}
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		overrides[configID] = value
	}
	if len(overrides) == 0 {
		overrides = nil
	}

	return strings.TrimSpace(raw.ModelID), overrides
}

func threadSessionID(agentOptionsJSON string) string {
	var raw struct {
		SessionID string `json:"sessionId"`
	}
	if strings.TrimSpace(agentOptionsJSON) == "" {
		return ""
	}
	if err := json.Unmarshal([]byte(agentOptionsJSON), &raw); err != nil {
		return ""
	}
	return strings.TrimSpace(raw.SessionID)
}

func threadFreshSessionRequested(agentOptionsJSON string) bool {
	var raw struct {
		FreshSession bool `json:"_ngentFreshSession"`
	}
	if strings.TrimSpace(agentOptionsJSON) == "" {
		return false
	}
	if err := json.Unmarshal([]byte(agentOptionsJSON), &raw); err != nil {
		return false
	}
	return raw.FreshSession
}

func withoutThreadSessionID(agentOptionsJSON string) (string, error) {
	objectValue := map[string]any{}
	trimmed := strings.TrimSpace(agentOptionsJSON)
	if trimmed != "" {
		if err := json.Unmarshal([]byte(trimmed), &objectValue); err != nil {
			return "", err
		}
	}

	delete(objectValue, "sessionId")
	delete(objectValue, threadAgentOptionFreshSessionKey)
	normalized, err := json.Marshal(objectValue)
	if err != nil {
		return "", err
	}
	return string(normalized), nil
}

func isSessionOnlyAgentOptionsUpdate(currentAgentOptionsJSON, nextAgentOptionsJSON string) (bool, error) {
	currentWithoutSessionID, err := withoutThreadSessionID(currentAgentOptionsJSON)
	if err != nil {
		return false, err
	}
	nextWithoutSessionID, err := withoutThreadSessionID(nextAgentOptionsJSON)
	if err != nil {
		return false, err
	}
	if currentWithoutSessionID == nextWithoutSessionID {
		return threadSessionID(currentAgentOptionsJSON) != threadSessionID(nextAgentOptionsJSON) ||
			threadFreshSessionRequested(currentAgentOptionsJSON) != threadFreshSessionRequested(nextAgentOptionsJSON), nil
	}

	currentComparable, err := withoutThreadConfigState(currentWithoutSessionID)
	if err != nil {
		return false, err
	}
	nextComparable, err := withoutThreadConfigState(nextWithoutSessionID)
	if err != nil {
		return false, err
	}
	if currentComparable != nextComparable {
		return false, nil
	}

	return threadSessionID(currentAgentOptionsJSON) != threadSessionID(nextAgentOptionsJSON) ||
		threadFreshSessionRequested(currentAgentOptionsJSON) != threadFreshSessionRequested(nextAgentOptionsJSON), nil
}

func withThreadSessionID(agentOptionsJSON, sessionID string) (string, bool, error) {
	sessionID = strings.TrimSpace(sessionID)
	previousNormalized := normalizeThreadAgentOptionsForScope(agentOptionsJSON)

	objectValue := map[string]any{}
	trimmed := strings.TrimSpace(agentOptionsJSON)
	if trimmed != "" {
		if err := json.Unmarshal([]byte(trimmed), &objectValue); err != nil {
			return "", false, err
		}
	}

	if sessionID == "" {
		delete(objectValue, "sessionId")
	} else {
		objectValue["sessionId"] = sessionID
		delete(objectValue, threadAgentOptionFreshSessionKey)
	}

	normalized, err := json.Marshal(objectValue)
	if err != nil {
		return "", false, err
	}
	nextNormalized := string(normalized)
	return nextNormalized, previousNormalized != nextNormalized, nil
}

func withThreadFreshSessionRequested(agentOptionsJSON string, fresh bool) (string, bool, error) {
	previousNormalized := normalizeThreadAgentOptionsForScope(agentOptionsJSON)

	objectValue := map[string]any{}
	trimmed := strings.TrimSpace(agentOptionsJSON)
	if trimmed != "" {
		if err := json.Unmarshal([]byte(trimmed), &objectValue); err != nil {
			return "", false, err
		}
	}

	if fresh {
		objectValue[threadAgentOptionFreshSessionKey] = true
	} else {
		delete(objectValue, threadAgentOptionFreshSessionKey)
	}

	normalized, err := json.Marshal(objectValue)
	if err != nil {
		return "", false, err
	}
	nextNormalized := string(normalized)
	return nextNormalized, previousNormalized != nextNormalized, nil
}

func sanitizeThreadAgentOptionsForResponse(agentOptionsJSON string) (json.RawMessage, error) {
	objectValue := map[string]any{}
	trimmed := strings.TrimSpace(agentOptionsJSON)
	if trimmed != "" {
		if err := json.Unmarshal([]byte(trimmed), &objectValue); err != nil {
			return nil, err
		}
	}

	delete(objectValue, threadAgentOptionFreshSessionKey)
	normalized, err := json.Marshal(objectValue)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(normalized), nil
}

func applyThreadConfigSelections(
	options []agents.ConfigOption,
	modelID string,
	overrides map[string]string,
) []agents.ConfigOption {
	cloned := acpmodel.CloneConfigOptions(options)
	modelID = strings.TrimSpace(modelID)

	for i := range cloned {
		configID := strings.TrimSpace(cloned[i].ID)
		if configID == "" {
			continue
		}
		if strings.EqualFold(configID, "model") || strings.EqualFold(strings.TrimSpace(cloned[i].Category), "model") {
			if modelID != "" {
				cloned[i].CurrentValue = modelID
			}
			continue
		}
		if len(overrides) == 0 {
			continue
		}
		if value := strings.TrimSpace(overrides[configID]); value != "" {
			cloned[i].CurrentValue = value
		}
	}

	return acpmodel.NormalizeConfigOptions(cloned)
}

func modelOnlyThreadConfigOptions(options []agents.ConfigOption) []agents.ConfigOption {
	modelOption, ok := acpmodel.FindModelConfigOption(options)
	if !ok {
		return nil
	}
	return acpmodel.NormalizeConfigOptions([]agents.ConfigOption{modelOption})
}
