package httpapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/observability"
	"github.com/beyond5959/ngent/internal/runtime"
	"github.com/beyond5959/ngent/internal/storage"
)

// AgentInfo describes one supported agent entry returned by /v1/agents.
type AgentInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// ThreadStore is the storage contract required by HTTP APIs.
type ThreadStore interface {
	UpsertClient(ctx context.Context, clientID string) error
	CreateThread(ctx context.Context, params storage.CreateThreadParams) (storage.Thread, error)
	GetThread(ctx context.Context, threadID string) (storage.Thread, error)
	DeleteThread(ctx context.Context, threadID string) error
	UpdateThreadTitle(ctx context.Context, threadID, title string) error
	UpdateThreadSummary(ctx context.Context, threadID, summary string) error
	UpdateThreadAgentOptions(ctx context.Context, threadID, agentOptionsJSON string) error
	UpsertAgentConfigCatalog(ctx context.Context, params storage.UpsertAgentConfigCatalogParams) error
	GetAgentConfigCatalog(ctx context.Context, agentID, modelID string) (storage.AgentConfigCatalog, error)
	ListAgentConfigCatalogsByAgent(ctx context.Context, agentID string) ([]storage.AgentConfigCatalog, error)
	GetAgentSlashCommands(ctx context.Context, agentID string) (storage.AgentSlashCommands, error)
	UpsertAgentSlashCommands(ctx context.Context, params storage.UpsertAgentSlashCommandsParams) error
	GetSessionTranscriptCache(ctx context.Context, agentID, cwd, sessionID string) (storage.SessionTranscriptCache, error)
	UpsertSessionTranscriptCache(ctx context.Context, params storage.UpsertSessionTranscriptCacheParams) error
	GetSessionConfigCache(ctx context.Context, agentID, cwd, sessionID string) (storage.SessionConfigCache, error)
	UpsertSessionConfigCache(ctx context.Context, params storage.UpsertSessionConfigCacheParams) error
	GetSessionUsageCache(ctx context.Context, agentID, cwd, sessionID string) (storage.SessionUsageCache, error)
	UpsertSessionUsageCache(ctx context.Context, params storage.UpsertSessionUsageCacheParams) error
	ListThreads(ctx context.Context) ([]storage.Thread, error)
	CreateTurn(ctx context.Context, params storage.CreateTurnParams) (storage.Turn, error)
	CreateTurnAttachments(ctx context.Context, params []storage.CreateTurnAttachmentParams) error
	GetTurnAttachment(ctx context.Context, attachmentID string) (storage.TurnAttachment, error)
	GetTurn(ctx context.Context, turnID string) (storage.Turn, error)
	ListTurnsByThread(ctx context.Context, threadID string) ([]storage.Turn, error)
	AppendEvent(ctx context.Context, turnID, eventType, dataJSON string) (storage.Event, error)
	ListEventsByTurns(ctx context.Context, turnIDs []string) (map[string][]storage.Event, error)
	ListEventsByTurn(ctx context.Context, turnID string) ([]storage.Event, error)
	FinalizeTurn(ctx context.Context, params storage.FinalizeTurnParams) error
	ListRecentDirectories(ctx context.Context, clientID string, limit int) ([]string, error)
}

// TurnAgentFactory resolves a per-turn agent provider from thread metadata.
type TurnAgentFactory func(thread storage.Thread) (agents.Streamer, error)

// AgentModelsFactory resolves selectable model options for one agent.
type AgentModelsFactory func(ctx context.Context, agentID string) ([]agents.ModelOption, error)

// Config controls HTTP API behavior.
type Config struct {
	AuthToken          string
	DataDir            string
	Agents             []AgentInfo
	AllowedAgentIDs    []string
	AllowedRoots       []string
	Store              ThreadStore
	TurnController     *runtime.TurnController
	Agent              agents.Streamer
	TurnAgentFactory   TurnAgentFactory
	AgentModelsFactory AgentModelsFactory
	AgentIdleTTL       time.Duration
	Logger             *observability.Logger
	ContextRecentTurns int
	ContextMaxChars    int
	CompactMaxChars    int
	PermissionTimeout  time.Duration
	FrontendHandler    http.Handler
}

// Server serves the HTTP API.
type Server struct {
	authToken          string
	dataDir            string
	agents             []AgentInfo
	allowedRoots       []string
	store              ThreadStore
	allowedAgent       map[string]struct{}
	turns              *runtime.TurnController
	turnAgentFactory   TurnAgentFactory
	agentModelsFactory AgentModelsFactory
	agentIdleTTL       time.Duration
	logger             *observability.Logger
	contextRecentTurns int
	contextMaxChars    int
	compactMaxChars    int
	permissionTimeout  time.Duration
	frontendHandler    http.Handler

	permissionsMu sync.Mutex
	permissions   map[string]*pendingPermission
	permissionSeq uint64

	agentMu       sync.Mutex
	agentsByScope map[string]*managedAgent
	turnEvents    *turnEventBroker
	janitorStop   chan struct{}
	janitorDone   chan struct{}
}

const (
	defaultContextRecentTurns = 10
	defaultContextMaxChars    = 20000
	defaultCompactMaxChars    = 4000
	defaultAgentIdleTTL       = 5 * time.Minute
	defaultPermissionTimeout  = 2 * time.Hour

	threadAgentOptionFreshSessionKey = "_ngentFreshSession"
	eventTypeUserPrompt              = "user_prompt"
	eventTypeMessageContent          = "message_content"
	eventTypePermissionResolved      = "permission_resolved"
	eventTypeReasoningDelta          = "reasoning_delta"
	eventTypeSessionInfoUpdate       = "session_info_update"
	eventTypeSessionUsageUpdate      = "session_usage_update"
	eventTypeToolCall                = "tool_call"
	eventTypeToolCallUpdate          = "tool_call_update"
)

const (
	codeInvalidArgument     = "INVALID_ARGUMENT"
	codeUnauthorized        = "UNAUTHORIZED"
	codeForbidden           = "FORBIDDEN"
	codeNotFound            = "NOT_FOUND"
	codeConflict            = "CONFLICT"
	codeTimeout             = "TIMEOUT"
	codeInternal            = "INTERNAL"
	codeUpstreamUnavailable = "UPSTREAM_UNAVAILABLE"
)

var errThreadConfigOptionsUnavailable = errors.New("thread config options are not available yet")

const maxTurnMultipartMemory = 32 << 20

type turnCreateRequest struct {
	Prompt     agents.Prompt
	Stream     bool
	FullAccess bool
	Uploads    []storedTurnAttachment
}

type storedTurnAttachment struct {
	PromptContent agents.PromptContent
	FilePath      string
}

// New creates a new API server.
func New(cfg Config) *Server {
	agentsList := make([]AgentInfo, len(cfg.Agents))
	copy(agentsList, cfg.Agents)

	roots := make([]string, 0, len(cfg.AllowedRoots))
	for _, root := range cfg.AllowedRoots {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}
		roots = append(roots, filepath.Clean(root))
	}

	allowedAgent := make(map[string]struct{}, len(cfg.AllowedAgentIDs))
	for _, agentID := range cfg.AllowedAgentIDs {
		agentID = strings.TrimSpace(agentID)
		if agentID == "" {
			continue
		}
		allowedAgent[agentID] = struct{}{}
	}

	turnController := cfg.TurnController
	if turnController == nil {
		turnController = runtime.NewTurnController()
	}

	turnAgentFactory := cfg.TurnAgentFactory
	if turnAgentFactory == nil {
		agent := cfg.Agent
		if agent == nil {
			agent = agents.NewFakeAgent()
		}
		turnAgentFactory = func(thread storage.Thread) (agents.Streamer, error) {
			_ = thread
			return agent, nil
		}
	}

	permissionTimeout := cfg.PermissionTimeout
	if permissionTimeout <= 0 {
		permissionTimeout = defaultPermissionTimeout
	}

	contextRecentTurns := cfg.ContextRecentTurns
	if contextRecentTurns <= 0 {
		contextRecentTurns = defaultContextRecentTurns
	}

	contextMaxChars := cfg.ContextMaxChars
	if contextMaxChars <= 0 {
		contextMaxChars = defaultContextMaxChars
	}

	compactMaxChars := cfg.CompactMaxChars
	if compactMaxChars <= 0 {
		compactMaxChars = defaultCompactMaxChars
	}

	agentIdleTTL := cfg.AgentIdleTTL
	if agentIdleTTL <= 0 {
		agentIdleTTL = defaultAgentIdleTTL
	}

	logger := cfg.Logger
	if logger == nil {
		logger = observability.NewLoggerWithWriter(io.Discard, observability.LevelError)
	}

	dataDir := filepath.Clean(strings.TrimSpace(cfg.DataDir))
	if dataDir == "." || dataDir == "" {
		dataDir = uploadTempDir()
	}

	server := &Server{
		authToken:          cfg.AuthToken,
		dataDir:            dataDir,
		agents:             agentsList,
		allowedRoots:       roots,
		store:              cfg.Store,
		allowedAgent:       allowedAgent,
		turns:              turnController,
		turnAgentFactory:   turnAgentFactory,
		agentModelsFactory: cfg.AgentModelsFactory,
		agentIdleTTL:       agentIdleTTL,
		logger:             logger,
		contextRecentTurns: contextRecentTurns,
		contextMaxChars:    contextMaxChars,
		compactMaxChars:    compactMaxChars,
		permissionTimeout:  permissionTimeout,
		frontendHandler:    cfg.FrontendHandler,
		permissions:        make(map[string]*pendingPermission),
		agentsByScope:      make(map[string]*managedAgent),
		turnEvents:         newTurnEventBroker(),
		janitorStop:        make(chan struct{}),
		janitorDone:        make(chan struct{}),
	}
	go server.idleJanitorLoop()
	return server
}

// ServeHTTP handles all HTTP requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startedAt := time.Now()
	loggingWriter := newLoggingResponseWriter(w)
	s.serveHTTP(loggingWriter, r)
	s.logRequestCompletion(r, loggingWriter, startedAt)
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if attachmentID, ok := parseAttachmentPath(r.URL.Path); ok {
		s.handleAttachment(w, r, attachmentID)
		return
	}

	if r.URL.Path == "/healthz" {
		s.handleHealthz(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/v1/") {
		if !s.isAuthorized(r) {
			writeError(w, http.StatusUnauthorized, codeUnauthorized, "missing or invalid bearer token", map[string]any{
				"header": "Authorization",
			})
			return
		}

		clientID := strings.TrimSpace(r.Header.Get("X-Client-ID"))
		if clientID == "" {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "missing required header X-Client-ID", map[string]any{
				"header": "X-Client-ID",
			})
			return
		}

		if s.store == nil {
			writeError(w, http.StatusInternalServerError, codeInternal, "storage is not configured", map[string]any{})
			return
		}

		if err := s.store.UpsertClient(r.Context(), clientID); err != nil {
			writeError(w, http.StatusInternalServerError, codeInternal, "failed to upsert client", map[string]any{
				"reason": err.Error(),
			})
			return
		}

		s.routeV1(w, r, clientID)
		return
	}

	if s.frontendHandler != nil {
		s.frontendHandler.ServeHTTP(w, r)
		return
	}

	writeError(w, http.StatusNotFound, codeNotFound, "endpoint not found", map[string]any{"path": r.URL.Path})
}

func (s *Server) logRequestCompletion(r *http.Request, w *loggingResponseWriter, startedAt time.Time) {
	if s.logger == nil {
		return
	}
	s.logger.HTTPRequest(observability.HTTPRequestLogEntry{
		RemoteAddr:  requestClientAddr(r),
		Method:      r.Method,
		Path:        requestLogPath(r),
		Proto:       r.Proto,
		Status:      w.StatusCode(),
		RequestTime: startedAt,
		Duration:    time.Since(startedAt),
	})
}

func (s *Server) routeV1(w http.ResponseWriter, r *http.Request, clientID string) {
	if r.URL.Path == "/v1/agents" {
		s.handleAgents(w, r)
		return
	}
	if agentID, ok := parseAgentModelsPath(r.URL.Path); ok {
		s.handleAgentModels(w, r, agentID)
		return
	}

	if r.URL.Path == "/v1/path-search" {
		s.handlePathSearch(w, r)
		return
	}

	if r.URL.Path == "/v1/recent-directories" {
		s.handleRecentDirectories(w, r, clientID)
		return
	}

	if r.URL.Path == "/v1/threads" {
		s.handleThreadsCollection(w, r, clientID)
		return
	}

	if permissionID, ok := parsePermissionPath(r.URL.Path); ok {
		s.handlePermissionDecision(w, r, clientID, permissionID)
		return
	}

	if turnID, ok := parseTurnEventsPath(r.URL.Path); ok {
		s.handleTurnEventsStream(w, r, clientID, turnID)
		return
	}

	if turnID, ok := parseTurnCancelPath(r.URL.Path); ok {
		s.handleCancelTurn(w, r, clientID, turnID)
		return
	}

	if threadID, subresource, ok := parseThreadPath(r.URL.Path); ok {
		s.handleThreadResource(w, r, clientID, threadID, subresource)
		return
	}

	writeError(w, http.StatusNotFound, codeNotFound, "endpoint not found", map[string]any{"path": r.URL.Path})
}
