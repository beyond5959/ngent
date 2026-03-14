package agents

import "context"

// ReasoningHandler receives one hidden reasoning delta for the active turn.
type ReasoningHandler func(ctx context.Context, delta string) error

type reasoningHandlerContextKey struct{}

// WithReasoningHandler binds one per-turn reasoning callback to context.
func WithReasoningHandler(ctx context.Context, handler ReasoningHandler) context.Context {
	if handler == nil {
		return ctx
	}
	return context.WithValue(ctx, reasoningHandlerContextKey{}, handler)
}

// ReasoningHandlerFromContext gets reasoning callback from context, if present.
func ReasoningHandlerFromContext(ctx context.Context) (ReasoningHandler, bool) {
	if ctx == nil {
		return nil, false
	}
	handler, ok := ctx.Value(reasoningHandlerContextKey{}).(ReasoningHandler)
	if !ok || handler == nil {
		return nil, false
	}
	return handler, true
}

// NotifyReasoningDelta reports one hidden reasoning delta to the active callback.
func NotifyReasoningDelta(ctx context.Context, delta string) error {
	handler, ok := ReasoningHandlerFromContext(ctx)
	if !ok || delta == "" {
		return nil
	}
	return handler(ctx, delta)
}
