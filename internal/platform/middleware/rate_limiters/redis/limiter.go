package redis

import "context"

type Limiter interface {
	Allow(ctx context.Context, key string) (bool, int, error)
	Name() string
}
