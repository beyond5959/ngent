package agents

import (
	"context"
	"strings"
)

// PromptRuntimeOverrides carries one-shot runtime overrides for a single turn.
type PromptRuntimeOverrides struct {
	ApprovalPolicy string
	Sandbox        string
}

type promptRuntimeOverridesContextKey struct{}

// WithPromptRuntimeOverrides binds one-shot runtime overrides to the turn context.
func WithPromptRuntimeOverrides(ctx context.Context, overrides PromptRuntimeOverrides) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	overrides = PromptRuntimeOverrides{
		ApprovalPolicy: strings.TrimSpace(overrides.ApprovalPolicy),
		Sandbox:        strings.TrimSpace(overrides.Sandbox),
	}
	if overrides.ApprovalPolicy == "" && overrides.Sandbox == "" {
		return ctx
	}
	return context.WithValue(ctx, promptRuntimeOverridesContextKey{}, overrides)
}

// PromptRuntimeOverridesFromContext gets one-shot runtime overrides from context, if present.
func PromptRuntimeOverridesFromContext(ctx context.Context) (PromptRuntimeOverrides, bool) {
	if ctx == nil {
		return PromptRuntimeOverrides{}, false
	}
	overrides, ok := ctx.Value(promptRuntimeOverridesContextKey{}).(PromptRuntimeOverrides)
	if !ok {
		return PromptRuntimeOverrides{}, false
	}
	if overrides.ApprovalPolicy == "" && overrides.Sandbox == "" {
		return PromptRuntimeOverrides{}, false
	}
	return overrides, true
}
