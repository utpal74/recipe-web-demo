package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type fixedWindowLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

// Name implements [Limiter].
func (f *fixedWindowLimiter) Name() string {
	return "fixed"
}

func NewFixedWindowLimiter(client *redis.Client, limit int, window time.Duration) Limiter {
	return &fixedWindowLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (f *fixedWindowLimiter) Allow(ctx context.Context, key string) (bool, int, int, int64, error) {

	window := time.Now().Unix() / int64(f.window.Seconds())

	key = fmt.Sprintf("%s:%d", key, window)

	count, err := f.client.Incr(ctx, key).Result()
	if err != nil {
		return false, 0, 0, 0, err
	}

	if count == 1 {
		err := f.client.Expire(ctx, key, f.window).Err()
		if err != nil {
			return false, 0, 0, 0, err
		}
	}

	if count > int64(f.limit) {
		return false, 0, 0, 0, err
	}

	remaining := max(f.limit-int(count), 0)

	reset := (window + 1) * int64(f.window.Seconds())

	return true, f.limit, remaining, reset, nil
}
