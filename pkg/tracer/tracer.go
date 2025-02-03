package tracer

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const traceIDKey ctxKey = "traceID"

func WithTrace(ctx context.Context) context.Context {
	val := TraceID(ctx)
	if val == "" {
		return context.WithValue(ctx, traceIDKey, uuid.NewString())
	}
	return ctx
}

func TraceID(ctx context.Context) string {
	val, ok := ctx.Value(traceIDKey).(string)
	if ok {
		return val
	}
	return ""
}
