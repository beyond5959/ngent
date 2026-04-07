package pi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/beyond5959/acp-adapter/pkg/piacp"
	"github.com/beyond5959/ngent/internal/agents"
	"github.com/beyond5959/ngent/internal/agents/acpmodel"
	"github.com/beyond5959/ngent/internal/observability"
)

var _ agents.SessionTranscriptLoader = (*Client)(nil)

// LoadSessionTranscript replays one Pi session through ACP session/load.
func (c *Client) LoadSessionTranscript(
	ctx context.Context,
	req agents.SessionTranscriptRequest,
) (agents.SessionTranscriptResult, error) {
	if c == nil {
		return agents.SessionTranscriptResult{}, errors.New("pi: nil client")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	startCtx, cancel := context.WithTimeout(ctx, c.startTimeout)
	defer cancel()

	runtime, caps, err := c.startRuntime(startCtx)
	if err != nil {
		return agents.SessionTranscriptResult{}, err
	}
	defer runtime.Close()

	if !caps.CanLoad {
		return agents.SessionTranscriptResult{}, agents.ErrSessionLoadUnsupported
	}

	session, err := agents.FindSessionByID(startCtx, c, req.CWD, req.SessionID)
	if err != nil {
		return agents.SessionTranscriptResult{}, err
	}

	loadCWD := piSessionCWD(c, session.CWD)
	if strings.TrimSpace(loadCWD) == "" {
		loadCWD = piSessionCWD(c, req.CWD)
	}

	return c.collectSessionReplay(ctx, runtime, session.SessionID, loadCWD)
}

func (c *Client) collectSessionReplay(
	ctx context.Context,
	runtime *piacp.EmbeddedRuntime,
	sessionID string,
	loadCWD string,
) (agents.SessionTranscriptResult, error) {
	if runtime == nil {
		return agents.SessionTranscriptResult{}, errors.New("pi: embedded runtime is nil")
	}

	collector := agents.NewACPTranscriptCollector()
	updates, unsubscribe := runtime.SubscribeUpdates(256)
	defer unsubscribe()

	type loadResult struct {
		err    error
		result json.RawMessage
	}
	loadDone := make(chan loadResult, 1)
	go func() {
		rawResult, err := c.clientRequest(ctx, runtime, "session/load", map[string]any{
			"sessionId":  sessionID,
			"cwd":        loadCWD,
			"mcpServers": []any{},
		})
		if err != nil {
			loadDone <- loadResult{err: fmt.Errorf("pi: session/load failed: %w", err)}
			return
		}
		loadDone <- loadResult{result: rawResult.Result}
	}()

	for {
		select {
		case <-ctx.Done():
			return agents.SessionTranscriptResult{}, ctx.Err()
		case result := <-loadDone:
			if result.err != nil {
				return agents.SessionTranscriptResult{}, result.err
			}
			if err := drainPiReplayUpdates(ctx, collector, updates, c.Name()); err != nil {
				return agents.SessionTranscriptResult{}, err
			}
			replay := collector.Result()
			replay.ConfigOptions = acpmodel.ExtractConfigOptions(result.result)
			return agents.CloneSessionTranscriptResult(replay), nil
		case msg, ok := <-updates:
			if !ok {
				if ctx.Err() != nil {
					return agents.SessionTranscriptResult{}, ctx.Err()
				}
				return agents.SessionTranscriptResult{}, errors.New("pi: embedded updates channel closed")
			}
			if err := consumePiReplayUpdate(ctx, collector, msg, c.Name()); err != nil {
				return agents.SessionTranscriptResult{}, err
			}
		}
	}
}

func drainPiReplayUpdates(
	ctx context.Context,
	collector *agents.ACPTranscriptCollector,
	updates <-chan piacp.RPCMessage,
	component string,
) error {
	for {
		select {
		case msg, ok := <-updates:
			if !ok {
				return nil
			}
			if err := consumePiReplayUpdate(ctx, collector, msg, component); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func consumePiReplayUpdate(
	ctx context.Context,
	collector *agents.ACPTranscriptCollector,
	msg piacp.RPCMessage,
	component string,
) error {
	observability.LogACPMessage(component, "inbound", msg)
	if msg.Method != methodSessionUpdate || len(msg.Params) == 0 {
		return nil
	}
	update, err := agents.ParseACPUpdate(msg.Params)
	if err != nil {
		return err
	}
	if update.Type == agents.ACPUpdateTypeUsage && update.SessionUsage != nil {
		if err := agents.NotifySessionUsageUpdate(ctx, *update.SessionUsage); err != nil {
			return err
		}
	}
	collector.HandleUpdate(update)
	return nil
}
