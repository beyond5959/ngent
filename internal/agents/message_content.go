package agents

import (
	"context"
	"encoding/json"
	"strings"
)

// ACPMessageContent is one normalized non-text assistant message content block.
type ACPMessageContent struct {
	Content    json.RawMessage `json:"content,omitempty"`
	HasContent bool            `json:"-"`
}

// EventPayload returns one JSON-serializable SSE/history payload for the message content.
func (event ACPMessageContent) EventPayload(turnID string) map[string]any {
	payload := map[string]any{
		"turnId": strings.TrimSpace(turnID),
	}
	if event.HasContent {
		payload["content"] = cloneACPMessageContentJSON(event.Content)
	}
	return payload
}

// CloneACPMessageContent returns a deep copy of one message-content event.
func CloneACPMessageContent(event ACPMessageContent) ACPMessageContent {
	return ACPMessageContent{
		Content:    cloneACPMessageContentJSON(event.Content),
		HasContent: event.HasContent,
	}
}

func cloneACPMessageContentJSON(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	return append(json.RawMessage(nil), raw...)
}

// MessageContentHandler receives one non-text assistant content block for the active turn.
type MessageContentHandler func(ctx context.Context, event ACPMessageContent) error

type messageContentHandlerContextKey struct{}

// WithMessageContentHandler binds one per-turn assistant content callback to context.
func WithMessageContentHandler(ctx context.Context, handler MessageContentHandler) context.Context {
	if handler == nil {
		return ctx
	}
	return context.WithValue(ctx, messageContentHandlerContextKey{}, handler)
}

// MessageContentHandlerFromContext gets assistant content callback from context, if present.
func MessageContentHandlerFromContext(ctx context.Context) (MessageContentHandler, bool) {
	if ctx == nil {
		return nil, false
	}
	handler, ok := ctx.Value(messageContentHandlerContextKey{}).(MessageContentHandler)
	if !ok || handler == nil {
		return nil, false
	}
	return handler, true
}

// NotifyMessageContent reports one non-text assistant message content block to the active callback.
func NotifyMessageContent(ctx context.Context, event ACPMessageContent) error {
	if !event.HasContent {
		return nil
	}
	handler, ok := MessageContentHandlerFromContext(ctx)
	if !ok {
		return nil
	}
	return handler(ctx, CloneACPMessageContent(event))
}
