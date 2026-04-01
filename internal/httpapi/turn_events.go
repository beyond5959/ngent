package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/beyond5959/ngent/internal/sse"
	"github.com/beyond5959/ngent/internal/storage"
)

const turnEventSubscriberBuffer = 128

type liveTurnEvent struct {
	event   storage.Event
	payload map[string]any
}

type turnEventBroker struct {
	mu          sync.Mutex
	nextID      uint64
	subscribers map[string]map[uint64]chan liveTurnEvent
}

func newTurnEventBroker() *turnEventBroker {
	return &turnEventBroker{
		subscribers: make(map[string]map[uint64]chan liveTurnEvent),
	}
}

func (b *turnEventBroker) Subscribe(turnID string) (<-chan liveTurnEvent, func()) {
	turnID = strings.TrimSpace(turnID)
	ch := make(chan liveTurnEvent, turnEventSubscriberBuffer)
	if turnID == "" {
		close(ch)
		return ch, func() {}
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.nextID++
	byTurn := b.subscribers[turnID]
	if byTurn == nil {
		byTurn = make(map[uint64]chan liveTurnEvent)
		b.subscribers[turnID] = byTurn
	}
	subID := b.nextID
	byTurn[subID] = ch

	return ch, func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		byTurn := b.subscribers[turnID]
		if byTurn == nil {
			return
		}
		if existing, ok := byTurn[subID]; ok {
			delete(byTurn, subID)
			close(existing)
		}
		if len(byTurn) == 0 {
			delete(b.subscribers, turnID)
		}
	}
}

func (b *turnEventBroker) Publish(event storage.Event, payload map[string]any) {
	turnID := strings.TrimSpace(event.TurnID)
	if turnID == "" {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	byTurn := b.subscribers[turnID]
	if len(byTurn) == 0 {
		return
	}

	live := liveTurnEvent{
		event:   event,
		payload: cloneTurnEventPayload(payload),
	}
	for subID, ch := range byTurn {
		select {
		case ch <- live:
		default:
			close(ch)
			delete(byTurn, subID)
		}
	}
	if len(byTurn) == 0 {
		delete(b.subscribers, turnID)
	}
}

func (b *turnEventBroker) CloseTurn(turnID string) {
	turnID = strings.TrimSpace(turnID)
	if turnID == "" {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	byTurn := b.subscribers[turnID]
	for subID, ch := range byTurn {
		delete(byTurn, subID)
		close(ch)
	}
	delete(b.subscribers, turnID)
}

func cloneTurnEventPayload(payload map[string]any) map[string]any {
	if len(payload) == 0 {
		return map[string]any{}
	}
	cloned := make(map[string]any, len(payload))
	for key, value := range payload {
		cloned[key] = value
	}
	return cloned
}

func payloadWithSeq(payload map[string]any, seq int) map[string]any {
	withSeq := cloneTurnEventPayload(payload)
	withSeq["seq"] = seq
	return withSeq
}

func (s *Server) appendTurnEvent(
	ctx context.Context,
	turnID, eventType string,
	payload map[string]any,
) (storage.Event, map[string]any, error) {
	dataJSON, err := json.Marshal(payload)
	if err != nil {
		return storage.Event{}, nil, err
	}
	event, err := s.store.AppendEvent(ctx, turnID, eventType, string(dataJSON))
	if err != nil {
		return storage.Event{}, nil, err
	}
	streamPayload := payloadWithSeq(payload, event.Seq)
	s.turnEvents.Publish(event, streamPayload)
	return event, streamPayload, nil
}

func (s *Server) replayableTurnEventPayload(event storage.Event) map[string]any {
	payload := map[string]any{}
	raw := strings.TrimSpace(event.DataJSON)
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &payload)
	}
	payload["seq"] = event.Seq
	return payload
}

func (s *Server) handleTurnEventsStream(w http.ResponseWriter, r *http.Request, clientID, turnID string) {
	_ = clientID
	if err := requireMethod(r, http.MethodGet); err != nil {
		writeMethodNotAllowed(w, r)
		return
	}

	turn, err := s.store.GetTurn(r.Context(), turnID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, codeNotFound, "turn not found", map[string]any{})
			return
		}
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to load turn", map[string]any{"reason": err.Error()})
		return
	}

	if _, ok := s.getAccessibleThread(r.Context(), turn.ThreadID); !ok {
		writeError(w, http.StatusNotFound, codeNotFound, "turn not found", map[string]any{})
		return
	}

	afterSeq := 0
	if raw := strings.TrimSpace(r.URL.Query().Get("after")); raw != "" {
		value, convErr := strconv.Atoi(raw)
		if convErr != nil || value < 0 {
			writeError(w, http.StatusBadRequest, codeInvalidArgument, "after must be a non-negative integer", map[string]any{
				"field":  "after",
				"reason": raw,
			})
			return
		}
		afterSeq = value
	}

	streamWriter, err := sse.NewWriter(w)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "SSE is not supported by response writer", map[string]any{})
		return
	}

	liveEvents, unsubscribe := s.turnEvents.Subscribe(turnID)
	defer unsubscribe()

	events, err := s.store.ListEventsByTurn(r.Context(), turnID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, codeInternal, "failed to list turn events", map[string]any{"reason": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	streamWriter.Flush()

	lastSeq := afterSeq
	for _, event := range events {
		if event.Seq <= afterSeq {
			continue
		}
		if err := streamWriter.Event(event.Type, s.replayableTurnEventPayload(event)); err != nil {
			return
		}
		lastSeq = event.Seq
	}

	if !s.turns.IsTurnActive(turnID) {
		return
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case live, ok := <-liveEvents:
			if !ok {
				return
			}
			if live.event.Seq <= lastSeq {
				continue
			}
			if err := streamWriter.Event(live.event.Type, live.payload); err != nil {
				return
			}
			lastSeq = live.event.Seq
			if live.event.Type == "turn_completed" {
				return
			}
		}
	}
}

func parseTurnEventsPath(path string) (turnID string, ok bool) {
	const prefix = "/v1/turns/"
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}
	rest := strings.TrimPrefix(path, prefix)
	if !strings.HasSuffix(rest, "/events") {
		return "", false
	}
	turnID = strings.TrimSuffix(rest, "/events")
	turnID = strings.TrimSpace(strings.Trim(turnID, "/"))
	if turnID == "" || strings.Contains(turnID, "/") {
		return "", false
	}
	return turnID, true
}
