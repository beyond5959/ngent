package pi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/beyond5959/acp-adapter/pkg/piacp"
	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/agents/acpmodel"
	"github.com/beyond5959/ngent/internal/agents/acpsession"
	"github.com/beyond5959/ngent/internal/agents/agentutil"
	"github.com/beyond5959/ngent/internal/observability"
)

const (
	jsonRPCVersion = "2.0"

	methodInitialize             = "initialize"
	methodSessionNew             = "session/new"
	methodSessionPrompt          = "session/prompt"
	methodSessionCancel          = "session/cancel"
	methodSessionSetConfigOption = "session/set_config_option"
	methodSessionUpdate          = "session/update"
	methodSessionRequestApproval = "session/request_permission"
)

const (
	defaultStartTimeout   = 30 * time.Second
	defaultRequestTimeout = 15 * time.Second

	initialSlashCommandsWait = 250 * time.Millisecond
	postPromptDrainTimeout   = 250 * time.Millisecond
)

// Config configures one embedded Pi runtime provider instance.
type Config struct {
	Dir             string
	ModelID         string
	SessionID       string
	ConfigOverrides map[string]string
	Name            string
	RuntimeConfig   piacp.RuntimeConfig
	StartTimeout    time.Duration
	RequestTimeout  time.Duration
}

// Client streams turn output through one in-process Pi ACP runtime.
type Client struct {
	*agentutil.State

	name string

	runtimeConfig  piacp.RuntimeConfig
	startTimeout   time.Duration
	requestTimeout time.Duration

	initMu sync.Mutex
	mu     sync.Mutex
	closed bool

	runtime            *piacp.EmbeddedRuntime
	sessionID          string
	updateUnsub        func()
	configOptions      []agents.ConfigOption
	canLoadSession     bool
	slashCommands      []agents.SlashCommand
	slashCommandsKnown bool
	slashCommandsReady chan struct{}
	sessionUsageByID   map[string]agents.SessionUsageUpdate

	requestSeq uint64
}

var _ agents.Streamer = (*Client)(nil)
var _ agents.ConfigOptionManager = (*Client)(nil)
var _ agents.SessionLister = (*Client)(nil)
var _ agents.SlashCommandsProvider = (*Client)(nil)
var _ agents.SessionTranscriptLoader = (*Client)(nil)
var _ io.Closer = (*Client)(nil)

// DefaultRuntimeConfig returns the default embedded Pi runtime configuration.
func DefaultRuntimeConfig() piacp.RuntimeConfig {
	cfg := piacp.DefaultRuntimeConfig()
	if value := strings.TrimSpace(os.Getenv("PI_BIN")); value != "" {
		cfg.PiBin = value
	}
	if value := strings.TrimSpace(os.Getenv("PI_PROVIDER")); value != "" {
		cfg.DefaultProvider = value
	}
	if value := strings.TrimSpace(os.Getenv("PI_MODEL")); value != "" {
		cfg.DefaultModel = value
	}
	if value := strings.TrimSpace(os.Getenv("PI_SESSION_DIR")); value != "" {
		cfg.SessionDir = value
	}
	if disable, ok := parseEnvBool(os.Getenv("PI_DISABLE_GATE")); ok {
		cfg.EnableGate = !disable
	}
	return cfg
}

func parseEnvBool(raw string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "":
		return false, false
	case "1", "true", "yes", "on":
		return true, true
	case "0", "false", "no", "off":
		return false, true
	default:
		return false, false
	}
}

// Preflight checks whether Pi runtime prerequisites are available on the host.
func Preflight(cfg piacp.RuntimeConfig) error {
	bin := strings.TrimSpace(cfg.PiBin)
	if bin == "" {
		bin = strings.TrimSpace(DefaultRuntimeConfig().PiBin)
	}
	if bin == "" {
		return errors.New("pi: binary path is empty")
	}
	if _, err := exec.LookPath(bin); err != nil {
		return fmt.Errorf("pi: binary %q not found: %w", bin, err)
	}
	return nil
}

// New constructs one embedded Pi provider.
func New(cfg Config) (*Client, error) {
	runtimeCfg := cfg.RuntimeConfig
	if strings.TrimSpace(runtimeCfg.PiBin) == "" &&
		len(runtimeCfg.PiArgs) == 0 &&
		strings.TrimSpace(runtimeCfg.DefaultProvider) == "" &&
		strings.TrimSpace(runtimeCfg.DefaultModel) == "" &&
		strings.TrimSpace(runtimeCfg.SessionDir) == "" &&
		!runtimeCfg.EnableGate &&
		!runtimeCfg.TraceJSON &&
		strings.TrimSpace(runtimeCfg.TraceJSONFile) == "" &&
		strings.TrimSpace(runtimeCfg.LogLevel) == "" &&
		strings.TrimSpace(runtimeCfg.PatchApplyMode) == "" &&
		len(runtimeCfg.Profiles) == 0 &&
		strings.TrimSpace(runtimeCfg.DefaultProfile) == "" {
		runtimeCfg = DefaultRuntimeConfig()
	}
	if err := Preflight(runtimeCfg); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(cfg.Name)
	if name == "" {
		name = "pi-embedded"
	}

	startTimeout := cfg.StartTimeout
	if startTimeout <= 0 {
		startTimeout = defaultStartTimeout
	}
	requestTimeout := cfg.RequestTimeout
	if requestTimeout <= 0 {
		requestTimeout = defaultRequestTimeout
	}

	state, err := agentutil.NewState(agents.AgentIDPi, agentutil.Config{
		Dir:             cfg.Dir,
		ModelID:         cfg.ModelID,
		SessionID:       cfg.SessionID,
		ConfigOverrides: cfg.ConfigOverrides,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		State:          state,
		name:           name,
		runtimeConfig:  runtimeCfg,
		startTimeout:   startTimeout,
		requestTimeout: requestTimeout,
	}, nil
}

// Name returns provider name.
func (c *Client) Name() string {
	if c == nil || c.name == "" {
		return "pi-embedded"
	}
	return c.name
}

// ConfigOptions returns current ACP session config options.
func (c *Client) ConfigOptions(ctx context.Context) ([]agents.ConfigOption, error) {
	if c == nil {
		return nil, errors.New("pi: nil client")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if _, _, err := c.ensureInitialized(ctx); err != nil {
		return nil, fmt.Errorf("pi: initialize runtime: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return acpmodel.CloneConfigOptions(c.configOptions), nil
}

// SlashCommands returns the latest slash-command snapshot after runtime init.
func (c *Client) SlashCommands(ctx context.Context) ([]agents.SlashCommand, bool, error) {
	if c == nil {
		return nil, false, errors.New("pi: nil client")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if _, _, err := c.ensureInitialized(ctx); err != nil {
		return nil, false, fmt.Errorf("pi: initialize runtime: %w", err)
	}
	c.waitForInitialSlashCommands(ctx)

	c.mu.Lock()
	known := c.slashCommandsKnown
	commands := agents.CloneSlashCommands(c.slashCommands)
	c.mu.Unlock()
	return commands, known, nil
}

// ListSessions queries ACP session/list for the current cwd.
func (c *Client) ListSessions(ctx context.Context, req agents.SessionListRequest) (agents.SessionListResult, error) {
	if c == nil {
		return agents.SessionListResult{}, errors.New("pi: nil client")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	startCtx, cancel := context.WithTimeout(ctx, c.startTimeout)
	defer cancel()

	runtime, caps, err := c.startRuntime(startCtx)
	if err != nil {
		return agents.SessionListResult{}, err
	}
	defer runtime.Close()

	if !caps.CanList || !caps.CanLoad {
		return agents.SessionListResult{}, agents.ErrSessionListUnsupported
	}

	params := map[string]any{
		"cwd":        piSessionCWD(c, req.CWD),
		"mcpServers": []any{},
	}
	if cursor := strings.TrimSpace(req.Cursor); cursor != "" {
		params["cursor"] = cursor
	}

	result, err := c.clientRequest(startCtx, runtime, "session/list", params)
	if err != nil {
		return agents.SessionListResult{}, fmt.Errorf("pi: session/list: %w", err)
	}
	return acpsession.ParseSessionListResult(result.Result)
}

// SetConfigOption applies one ACP session config option and returns latest options.
func (c *Client) SetConfigOption(ctx context.Context, configID, value string) ([]agents.ConfigOption, error) {
	if c == nil {
		return nil, errors.New("pi: nil client")
	}
	configID = strings.TrimSpace(configID)
	value = strings.TrimSpace(value)
	if configID == "" {
		return nil, errors.New("pi: configID is required")
	}
	if value == "" {
		return nil, errors.New("pi: value is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	runtime, sessionID, err := c.ensureInitialized(ctx)
	if err != nil {
		return nil, fmt.Errorf("pi: initialize runtime: %w", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()
	resp, err := c.clientRequest(reqCtx, runtime, methodSessionSetConfigOption, map[string]any{
		"sessionId": sessionID,
		"configId":  configID,
		"value":     value,
	})
	if err != nil {
		return nil, fmt.Errorf("pi: session/set_config_option failed: %w", err)
	}

	options := acpmodel.ExtractConfigOptions(resp.Result)
	c.mu.Lock()
	c.configOptions = acpmodel.CloneConfigOptions(options)
	c.mu.Unlock()
	c.ApplyConfigOptionResult(configID, value, options)
	return acpmodel.CloneConfigOptions(options), nil
}

// Close closes the embedded runtime.
func (c *Client) Close() error {
	if c == nil {
		return nil
	}

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	runtime := c.runtime
	c.runtime = nil
	c.sessionID = ""
	updateUnsub := c.updateUnsub
	c.updateUnsub = nil
	c.configOptions = nil
	c.canLoadSession = false
	c.slashCommands = nil
	c.slashCommandsKnown = false
	c.sessionUsageByID = nil
	slashCommandsReady := c.slashCommandsReady
	c.slashCommandsReady = nil
	c.mu.Unlock()

	if updateUnsub != nil {
		updateUnsub()
	}
	closeReadySignal(slashCommandsReady)
	if runtime != nil {
		return runtime.Close()
	}
	return nil
}

// Stream sends one prompt to embedded runtime and emits deltas.
func (c *Client) Stream(ctx context.Context, input string, onDelta func(delta string) error) (agents.StopReason, error) {
	return c.StreamPrompt(ctx, agents.TextPrompt(input), onDelta)
}

// StreamPrompt sends one structured prompt to embedded runtime and emits deltas.
func (c *Client) StreamPrompt(ctx context.Context, prompt agents.Prompt, onDelta func(delta string) error) (agents.StopReason, error) {
	if c == nil {
		return agents.StopReasonEndTurn, errors.New("pi: nil client")
	}
	if onDelta == nil {
		return agents.StopReasonEndTurn, errors.New("pi: onDelta callback is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	prompt = agents.NormalizePrompt(prompt)

	runtime, sessionID, err := c.ensureInitialized(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return agents.StopReasonCancelled, nil
		}
		return agents.StopReasonEndTurn, fmt.Errorf("pi: initialize runtime: %w", err)
	}
	c.waitForInitialSlashCommands(ctx)
	if err := c.notifyCachedSlashCommands(ctx); err != nil {
		return agents.StopReasonEndTurn, fmt.Errorf("pi: report slash commands: %w", err)
	}

	c.mu.Lock()
	configOptions := acpmodel.CloneConfigOptions(c.configOptions)
	c.mu.Unlock()
	if err := agents.NotifyConfigOptions(ctx, configOptions); err != nil {
		return agents.StopReasonEndTurn, fmt.Errorf("pi: report config options: %w", err)
	}
	if c.supportsLoadSession() {
		if err := agents.NotifySessionBound(ctx, sessionID); err != nil {
			return agents.StopReasonEndTurn, fmt.Errorf("pi: report session bound: %w", err)
		}
	}
	if err := c.notifyCachedSessionUsage(ctx, sessionID); err != nil {
		return agents.StopReasonEndTurn, fmt.Errorf("pi: report cached session usage: %w", err)
	}

	return c.streamOnce(ctx, runtime, sessionID, prompt, onDelta)
}

func (c *Client) streamOnce(
	ctx context.Context,
	runtime *piacp.EmbeddedRuntime,
	sessionID string,
	prompt agents.Prompt,
	onDelta func(delta string) error,
) (agents.StopReason, error) {
	updates, unsubscribe := runtime.SubscribeUpdates(256)
	defer unsubscribe()

	promptCtx, promptCancel := context.WithCancel(ctx)
	defer promptCancel()

	var stopWatchOnce sync.Once
	stopWatch := make(chan struct{})
	stopCancelWatcher := func() {
		stopWatchOnce.Do(func() {
			close(stopWatch)
		})
	}
	defer stopCancelWatcher()

	go func() {
		select {
		case <-promptCtx.Done():
			c.sendSessionCancel(runtime, sessionID)
		case <-stopWatch:
		}
	}()

	type promptResult struct {
		response piacp.RPCMessage
		err      error
	}
	promptDone := make(chan promptResult, 1)
	promptContent := prompt.ACPContent()
	if promptContent == nil {
		promptContent = []map[string]any{}
	}
	go func() {
		resp, reqErr := c.clientRequest(promptCtx, runtime, methodSessionPrompt, map[string]any{
			"sessionId": sessionID,
			"prompt":    promptContent,
		})
		promptDone <- promptResult{response: resp, err: reqErr}
	}()

	var (
		finalStopReason agents.StopReason
		promptFinished  bool
		drainTimer      *time.Timer
		drainCh         <-chan time.Time
	)
	stopDrainTimer := func() {
		if drainTimer == nil {
			return
		}
		if !drainTimer.Stop() {
			select {
			case <-drainTimer.C:
			default:
			}
		}
		drainCh = nil
	}
	resetDrainTimer := func() {
		if drainTimer == nil {
			drainTimer = time.NewTimer(postPromptDrainTimeout)
			drainCh = drainTimer.C
			return
		}
		if !drainTimer.Stop() {
			select {
			case <-drainTimer.C:
			default:
			}
		}
		drainTimer.Reset(postPromptDrainTimeout)
		drainCh = drainTimer.C
	}
	defer stopDrainTimer()

	for {
		select {
		case <-ctx.Done():
			stopCancelWatcher()
			stopDrainTimer()
			return agents.StopReasonCancelled, nil
		case result := <-promptDone:
			if result.err != nil {
				stopCancelWatcher()
				stopDrainTimer()
				if errors.Is(result.err, context.Canceled) || errors.Is(result.err, context.DeadlineExceeded) || ctx.Err() != nil {
					return agents.StopReasonCancelled, nil
				}
				return agents.StopReasonEndTurn, fmt.Errorf("pi: session/prompt failed: %w", result.err)
			}

			stopReason, parseErr := parsePromptStopReason(result.response.Result)
			if parseErr != nil {
				stopCancelWatcher()
				stopDrainTimer()
				return agents.StopReasonEndTurn, parseErr
			}
			usageUpdate, usageErr := agents.ParseACPPromptUsage(result.response.Result)
			if usageErr != nil {
				stopCancelWatcher()
				stopDrainTimer()
				return agents.StopReasonEndTurn, fmt.Errorf("pi: decode session/prompt usage: %w", usageErr)
			}
			if usageUpdate.SessionID == "" {
				usageUpdate.SessionID = strings.TrimSpace(sessionID)
			}
			c.cacheSessionUsage(usageUpdate)
			if err := agents.NotifySessionUsageUpdate(ctx, usageUpdate); err != nil {
				stopCancelWatcher()
				stopDrainTimer()
				return agents.StopReasonEndTurn, fmt.Errorf("pi: report session usage: %w", err)
			}
			if stopReason == "cancelled" {
				finalStopReason = agents.StopReasonCancelled
			} else {
				finalStopReason = agents.StopReasonEndTurn
			}
			promptFinished = true
			resetDrainTimer()
		case msg, ok := <-updates:
			if !ok {
				stopCancelWatcher()
				stopDrainTimer()
				if promptFinished {
					return finalStopReason, nil
				}
				if ctx.Err() != nil {
					return agents.StopReasonCancelled, nil
				}
				return agents.StopReasonEndTurn, errors.New("pi: embedded updates channel closed")
			}

			if err := c.handleUpdate(ctx, runtime, msg, onDelta); err != nil {
				stopCancelWatcher()
				stopDrainTimer()
				return agents.StopReasonEndTurn, err
			}
			if promptFinished {
				if acpSessionUpdateIsTerminal(msg.Params) {
					stopCancelWatcher()
					stopDrainTimer()
					return finalStopReason, nil
				}
				resetDrainTimer()
			}
		case <-drainCh:
			stopCancelWatcher()
			stopDrainTimer()
			if promptFinished {
				return finalStopReason, nil
			}
		}
	}
}

func (c *Client) handleUpdate(
	ctx context.Context,
	runtime *piacp.EmbeddedRuntime,
	msg piacp.RPCMessage,
	onDelta func(delta string) error,
) error {
	observability.LogACPMessage(c.Name(), "inbound", msg)

	if msg.Method == methodSessionUpdate {
		updateType := acpSessionUpdateTopLevelType(msg.Params)
		update, err := agents.ParseACPUpdate(msg.Params)
		if err != nil {
			return fmt.Errorf("pi: %w", err)
		}
		switch update.Type {
		case agents.ACPUpdateTypeMessageChunk:
			if update.Delta != "" {
				if err := onDelta(update.Delta); err != nil {
					c.sendSessionCancel(runtime, c.currentSessionID())
					return err
				}
				return nil
			}
			if update.MessageContent != nil {
				if err := agents.NotifyMessageContent(ctx, *update.MessageContent); err != nil {
					c.sendSessionCancel(runtime, c.currentSessionID())
					return err
				}
			}
			return nil
		case agents.ACPUpdateTypeThoughtMessageChunk:
			if updateType != "" && updateType != "reasoning" {
				return nil
			}
			if err := agents.NotifyReasoningDelta(ctx, update.Delta); err != nil {
				c.sendSessionCancel(runtime, c.currentSessionID())
				return err
			}
			return nil
		case agents.ACPUpdateTypePlan:
			handler, ok := agents.PlanHandlerFromContext(ctx)
			if !ok {
				return nil
			}
			if err := handler(ctx, update.PlanEntries); err != nil {
				c.sendSessionCancel(runtime, c.currentSessionID())
				return err
			}
		case agents.ACPUpdateTypeSessionInfo:
			if update.SessionInfo == nil {
				return nil
			}
			if err := agents.NotifySessionInfoUpdate(ctx, *update.SessionInfo); err != nil {
				c.sendSessionCancel(runtime, c.currentSessionID())
				return err
			}
		case agents.ACPUpdateTypeUsage:
			if update.SessionUsage == nil {
				return nil
			}
			c.cacheSessionUsage(*update.SessionUsage)
			if err := agents.NotifySessionUsageUpdate(ctx, *update.SessionUsage); err != nil {
				c.sendSessionCancel(runtime, c.currentSessionID())
				return err
			}
		case agents.ACPUpdateTypeAvailableCommands:
			c.cacheSlashCommands(update.Commands)
			if err := agents.NotifySlashCommands(ctx, update.Commands); err != nil {
				c.sendSessionCancel(runtime, c.currentSessionID())
				return err
			}
		case agents.ACPUpdateTypeToolCall, agents.ACPUpdateTypeToolCallUpdate:
			if update.ToolCall == nil {
				return nil
			}
			if err := agents.NotifyToolCall(ctx, *update.ToolCall); err != nil {
				c.sendSessionCancel(runtime, c.currentSessionID())
				return err
			}
		}
		return nil
	}

	if msg.Method == methodSessionRequestApproval {
		return c.handlePermissionRequest(ctx, runtime, msg)
	}

	if msg.Method != "" && msg.ID != nil {
		c.sendSessionCancel(runtime, c.currentSessionID())
		return fmt.Errorf("pi: unsupported embedded request method %q", msg.Method)
	}
	return nil
}

func acpSessionUpdateTopLevelType(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var payload struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}
	return strings.TrimSpace(payload.Type)
}

func acpSessionUpdateIsTerminal(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var payload struct {
		Type   string `json:"type"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return false
	}
	if strings.TrimSpace(payload.Type) != "status" {
		return false
	}
	switch strings.TrimSpace(payload.Status) {
	case "turn_completed", "turn_cancelled":
		return true
	default:
		return false
	}
}

func (c *Client) installUpdateMonitor(runtime *piacp.EmbeddedRuntime) {
	if runtime == nil {
		return
	}

	updates, unsubscribe := runtime.SubscribeUpdates(256)
	ready := make(chan struct{})

	c.mu.Lock()
	prevUnsub := c.updateUnsub
	prevReady := c.slashCommandsReady
	c.updateUnsub = unsubscribe
	c.slashCommands = nil
	c.slashCommandsKnown = false
	c.slashCommandsReady = ready
	c.sessionUsageByID = nil
	c.mu.Unlock()

	if prevUnsub != nil {
		prevUnsub()
	}
	closeReadySignal(prevReady)

	go c.monitorUpdates(updates, ready)
}

func (c *Client) monitorUpdates(updates <-chan piacp.RPCMessage, ready chan struct{}) {
	defer closeReadySignal(ready)

	for msg := range updates {
		if msg.Method != methodSessionUpdate || len(msg.Params) == 0 {
			continue
		}
		update, err := agents.ParseACPUpdate(msg.Params)
		if err != nil {
			continue
		}
		if update.Type == agents.ACPUpdateTypeUsage && update.SessionUsage != nil {
			c.cacheSessionUsage(*update.SessionUsage)
		}
		if update.Type != agents.ACPUpdateTypeAvailableCommands {
			continue
		}
		c.mu.Lock()
		if c.slashCommandsReady == ready {
			c.slashCommands = agents.CloneSlashCommands(update.Commands)
			c.slashCommandsKnown = true
		}
		c.mu.Unlock()
		closeReadySignal(ready)
	}
}

func (c *Client) clearUpdateMonitor() {
	c.mu.Lock()
	updateUnsub := c.updateUnsub
	c.updateUnsub = nil
	slashCommandsReady := c.slashCommandsReady
	c.slashCommandsReady = nil
	c.slashCommands = nil
	c.slashCommandsKnown = false
	c.sessionUsageByID = nil
	c.mu.Unlock()

	if updateUnsub != nil {
		updateUnsub()
	}
	closeReadySignal(slashCommandsReady)
}

func (c *Client) cacheSessionUsage(update agents.SessionUsageUpdate) {
	update = agents.CloneSessionUsageUpdate(update)
	if update.SessionID == "" || !agents.HasSessionUsageValues(update) {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sessionUsageByID == nil {
		c.sessionUsageByID = make(map[string]agents.SessionUsageUpdate)
	}
	c.sessionUsageByID[update.SessionID] = agents.MergeSessionUsageUpdate(c.sessionUsageByID[update.SessionID], update)
}

func (c *Client) notifyCachedSessionUsage(ctx context.Context, sessionID string) error {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil
	}

	c.mu.Lock()
	update, ok := c.sessionUsageByID[sessionID]
	c.mu.Unlock()
	if !ok {
		return nil
	}
	return agents.NotifySessionUsageUpdate(ctx, agents.CloneSessionUsageUpdate(update))
}

func (c *Client) cacheSlashCommands(commands []agents.SlashCommand) {
	c.mu.Lock()
	c.slashCommands = agents.CloneSlashCommands(commands)
	c.slashCommandsKnown = true
	c.mu.Unlock()
}

func (c *Client) waitForInitialSlashCommands(ctx context.Context) {
	c.mu.Lock()
	if c.slashCommandsKnown || c.slashCommandsReady == nil {
		c.mu.Unlock()
		return
	}
	ready := c.slashCommandsReady
	c.mu.Unlock()

	timer := time.NewTimer(initialSlashCommandsWait)
	defer timer.Stop()

	select {
	case <-ready:
	case <-ctx.Done():
	case <-timer.C:
	}
}

func (c *Client) notifyCachedSlashCommands(ctx context.Context) error {
	c.mu.Lock()
	known := c.slashCommandsKnown
	commands := agents.CloneSlashCommands(c.slashCommands)
	c.mu.Unlock()
	if !known {
		return nil
	}
	return agents.NotifySlashCommands(ctx, commands)
}

func (c *Client) handlePermissionRequest(
	ctx context.Context,
	runtime *piacp.EmbeddedRuntime,
	msg piacp.RPCMessage,
) error {
	if msg.ID == nil {
		return errors.New("pi: permission request missing id")
	}

	rawParams := map[string]any{}
	if len(msg.Params) > 0 {
		if err := json.Unmarshal(msg.Params, &rawParams); err != nil {
			return fmt.Errorf("pi: decode permission params: %w", err)
		}
	}

	request := agents.PermissionRequest{
		RequestID: idToString(*msg.ID),
		Approval:  mapString(rawParams, "approval"),
		Command:   mapString(rawParams, "command"),
		RawParams: rawParams,
	}

	outcome := agents.PermissionOutcomeDeclined
	if handler, ok := agents.PermissionHandlerFromContext(ctx); ok {
		resp, err := handler(ctx, request)
		if err == nil {
			switch resp.Outcome {
			case agents.PermissionOutcomeApproved, agents.PermissionOutcomeDeclined, agents.PermissionOutcomeCancelled:
				outcome = resp.Outcome
			}
		}
	}

	respondCtx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	if err := runtime.RespondPermission(
		respondCtx,
		*msg.ID,
		piacp.PermissionDecision{Outcome: string(outcome)},
	); err != nil {
		return fmt.Errorf("pi: respond permission outcome: %w", err)
	}
	observability.LogACPMessage(c.Name(), "outbound", map[string]any{
		"jsonrpc": jsonRPCVersion,
		"id":      *msg.ID,
		"result": map[string]any{
			"outcome": string(outcome),
		},
	})
	return nil
}

func (c *Client) sendSessionCancel(runtime *piacp.EmbeddedRuntime, sessionID string) {
	if runtime == nil || strings.TrimSpace(sessionID) == "" {
		return
	}
	cancelCtx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	_, _ = c.clientRequest(cancelCtx, runtime, methodSessionCancel, map[string]any{
		"sessionId": sessionID,
	})
}

func (c *Client) currentSessionID() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sessionID
}

func (c *Client) supportsLoadSession() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.canLoadSession
}

func (c *Client) ensureInitialized(ctx context.Context) (*piacp.EmbeddedRuntime, string, error) {
	c.initMu.Lock()
	defer c.initMu.Unlock()

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil, "", errors.New("pi: client is closed")
	}
	if c.runtime != nil && c.sessionID != "" {
		runtime := c.runtime
		sessionID := c.sessionID
		c.mu.Unlock()
		return runtime, sessionID, nil
	}
	c.mu.Unlock()

	startCtx, cancel := context.WithTimeout(ctx, c.startTimeout)
	defer cancel()

	runtime, caps, err := c.startRuntime(startCtx)
	if err != nil {
		return nil, "", err
	}
	c.installUpdateMonitor(runtime)

	sessionID := c.CurrentSessionID()
	configOptions := []agents.ConfigOption(nil)
	if sessionID != "" {
		if !caps.CanLoad {
			c.clearUpdateMonitor()
			_ = runtime.Close()
			return nil, "", agents.ErrSessionLoadUnsupported
		}
		loadResp, err := c.clientRequest(startCtx, runtime, "session/load", map[string]any{
			"sessionId":  sessionID,
			"cwd":        c.Dir(),
			"mcpServers": []any{},
		})
		if err != nil {
			c.clearUpdateMonitor()
			_ = runtime.Close()
			return nil, "", fmt.Errorf("pi: session/load failed: %w", err)
		}
		configOptions = acpmodel.ExtractConfigOptions(loadResp.Result)
	} else {
		newParams := map[string]any{
			"cwd": c.Dir(),
		}
		if modelID := c.CurrentModelID(); modelID != "" {
			newParams["model"] = modelID
		}
		sessionResp, err := c.clientRequest(startCtx, runtime, methodSessionNew, newParams)
		if err != nil {
			c.clearUpdateMonitor()
			_ = runtime.Close()
			return nil, "", err
		}

		sessionID, err = parseSessionID(sessionResp.Result)
		if err != nil {
			c.clearUpdateMonitor()
			_ = runtime.Close()
			return nil, "", err
		}
		configOptions = acpmodel.ExtractConfigOptions(sessionResp.Result)
	}
	configOptions, err = c.applySessionSelections(startCtx, runtime, sessionID, configOptions)
	if err != nil {
		c.clearUpdateMonitor()
		_ = runtime.Close()
		return nil, "", err
	}
	c.ApplyConfigOptionsSnapshot(configOptions)

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		c.clearUpdateMonitor()
		_ = runtime.Close()
		return nil, "", errors.New("pi: client is closed")
	}

	c.runtime = runtime
	c.sessionID = sessionID
	c.configOptions = acpmodel.CloneConfigOptions(configOptions)
	c.canLoadSession = caps.CanLoad
	c.mu.Unlock()
	return runtime, sessionID, nil
}

func (c *Client) applySessionSelections(
	ctx context.Context,
	runtime *piacp.EmbeddedRuntime,
	sessionID string,
	options []agents.ConfigOption,
) ([]agents.ConfigOption, error) {
	current := options

	if modelID := strings.TrimSpace(c.CurrentModelID()); modelID != "" &&
		strings.TrimSpace(acpmodel.CurrentValueForConfig(current, "model")) != modelID {
		resp, err := c.clientRequest(ctx, runtime, methodSessionSetConfigOption, map[string]any{
			"sessionId": sessionID,
			"configId":  "model",
			"value":     modelID,
		})
		if err != nil {
			return nil, fmt.Errorf("pi: session/set_config_option(model) failed: %w", err)
		}
		if updated := acpmodel.ExtractConfigOptions(resp.Result); len(updated) > 0 {
			current = updated
		}
	}

	overrides := c.CurrentConfigOverrides()
	if len(overrides) == 0 {
		return current, nil
	}

	configIDs := make([]string, 0, len(overrides))
	for configID := range overrides {
		configIDs = append(configIDs, configID)
	}
	sort.Strings(configIDs)

	for _, configID := range configIDs {
		value := strings.TrimSpace(overrides[configID])
		if value == "" {
			continue
		}
		resp, err := c.clientRequest(ctx, runtime, methodSessionSetConfigOption, map[string]any{
			"sessionId": sessionID,
			"configId":  configID,
			"value":     value,
		})
		if err != nil {
			return nil, fmt.Errorf("pi: session/set_config_option(%s) failed: %w", configID, err)
		}
		if updated := acpmodel.ExtractConfigOptions(resp.Result); len(updated) > 0 {
			current = updated
		}
	}
	return current, nil
}

func (c *Client) startRuntime(
	ctx context.Context,
) (*piacp.EmbeddedRuntime, acpsession.Capabilities, error) {
	runtime := piacp.NewEmbeddedRuntime(c.runtimeConfig)
	if err := runtime.Start(context.Background()); err != nil {
		_ = runtime.Close()
		return nil, acpsession.Capabilities{}, err
	}

	initResp, err := c.clientRequest(ctx, runtime, methodInitialize, map[string]any{
		"client": map[string]any{
			"name": "ngent",
		},
	})
	if err != nil {
		_ = runtime.Close()
		return nil, acpsession.Capabilities{}, err
	}
	return runtime, acpsession.ParseInitializeCapabilities(initResp.Result), nil
}

func piSessionCWD(c *Client, cwd string) string {
	cwd = strings.TrimSpace(cwd)
	if cwd != "" {
		return cwd
	}
	return c.Dir()
}

func (c *Client) clientRequest(
	ctx context.Context,
	runtime *piacp.EmbeddedRuntime,
	method string,
	params any,
) (piacp.RPCMessage, error) {
	if runtime == nil {
		return piacp.RPCMessage{}, errors.New("pi: embedded runtime is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.nextRequestID()
	msg := piacp.RPCMessage{
		JSONRPC: jsonRPCVersion,
		ID:      &id,
		Method:  method,
	}

	if params != nil {
		paramsJSON, err := json.Marshal(params)
		if err != nil {
			return piacp.RPCMessage{}, fmt.Errorf("pi: marshal %s params: %w", method, err)
		}
		msg.Params = paramsJSON
	}
	observability.LogACPMessage(c.Name(), "outbound", msg)

	response, err := runtime.ClientRequest(ctx, msg)
	if err != nil {
		return piacp.RPCMessage{}, err
	}
	observability.LogACPMessage(c.Name(), "inbound", response)
	if response.Error != nil {
		return piacp.RPCMessage{}, fmt.Errorf(
			"pi: %s rpc error code=%d message=%s",
			method,
			response.Error.Code,
			strings.TrimSpace(response.Error.Message),
		)
	}
	return response, nil
}

func (c *Client) nextRequestID() json.RawMessage {
	id := atomic.AddUint64(&c.requestSeq, 1)
	raw := strconv.Quote(fmt.Sprintf("srv-%d", id))
	return json.RawMessage(raw)
}

func parseSessionID(raw json.RawMessage) (string, error) {
	var payload struct {
		SessionID string `json:"sessionId"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return "", fmt.Errorf("pi: decode session/new result: %w", err)
	}
	sessionID := strings.TrimSpace(payload.SessionID)
	if sessionID == "" {
		return "", errors.New("pi: session/new returned empty sessionId")
	}
	return sessionID, nil
}

func parsePromptStopReason(raw json.RawMessage) (string, error) {
	var payload struct {
		StopReason string `json:"stopReason"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return "", fmt.Errorf("pi: decode session/prompt result: %w", err)
	}
	stopReason := strings.TrimSpace(payload.StopReason)
	if stopReason == "" {
		stopReason = string(agents.StopReasonEndTurn)
	}
	return stopReason, nil
}

func mapString(values map[string]any, key string) string {
	value, _ := values[key]
	text, _ := value.(string)
	return text
}

func idToString(raw json.RawMessage) string {
	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		return asString
	}
	var asNumber float64
	if err := json.Unmarshal(raw, &asNumber); err == nil {
		return strconv.FormatFloat(asNumber, 'f', -1, 64)
	}
	return strings.TrimSpace(string(raw))
}

func closeReadySignal(ch chan struct{}) {
	if ch == nil {
		return
	}
	select {
	case <-ch:
		return
	default:
		close(ch)
	}
}
