package core

import "context"

// context key is unexported and so no other package can create it accidentally
type contextKey struct{}

// withRequestID returns a copy of ctx that carries the request ID.
func withRequestID(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, contextKey{}, requestId)
}

func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}
