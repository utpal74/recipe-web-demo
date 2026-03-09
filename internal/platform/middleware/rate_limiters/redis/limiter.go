package redis

import "context"

type Limiter interface {
	Allow(ctx context.Context, key string) (bool, int, int, int64, error)
	Name() string
}
