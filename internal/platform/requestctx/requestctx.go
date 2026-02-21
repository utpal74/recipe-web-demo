package requestctx

import (
	"context"

	"github.com/gin-demo/recipes-web/internal/module/auth/domain"
)

type contextKey int

const (
	contextKeyAuthIdentity contextKey = iota
	contextKeyIdempotency
)

func SetAuthIdentity(ctx context.Context, identity domain.Identity) context.Context {
	return context.WithValue(ctx, contextKeyAuthIdentity, identity)
}

func GetAuthIdentity(ctx context.Context) (domain.Identity, bool) {
	v, ok := ctx.Value(contextKeyAuthIdentity).(domain.Identity)
	return v, ok
}

func SetIdempotencyKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, contextKeyIdempotency, key)
}

func GetIdempotencyKey(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeyIdempotency).(string)
	return v, ok
}
