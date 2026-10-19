package core

import "context"

type contextKey struct{}

func withRequestID(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, contextKey{}, requestId)
}

func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}
