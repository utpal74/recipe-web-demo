package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var script = redis.NewScript(`
local key = KEYS[1]

local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local request_id = ARGV[4]

local window_start = now - window

redis.call("ZREMRANGEBYSCORE", key, "-inf", window_start)

local count = redis.call("ZCARD", key)

local oldest = redis.call("ZRANGE", key, 0, 0, "WITHSCORES")

if count >= limit then
   return {0, count, oldest}
end

redis.call("ZADD", key, now, request_id)
redis.call("PEXPIRE", key, window)

oldest = redis.call("ZRANGE", key, 0, 0, "WITHSCORES")

return {1, count+1, oldest}
	`)

type slideWindowLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

// Name implements [Limiter].
func (s *slideWindowLimiter) Name() string {
	return "slide"
}

func NewSlideWindowLimiter(client *redis.Client, limit int, window time.Duration) Limiter {
	return &slideWindowLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (s *slideWindowLimiter) Allow(ctx context.Context, key string) (bool, int, int, int64, error) {
	window := s.window.Milliseconds()
	now := time.Now().UnixMilli()
	requestID := uuid.NewString()

	res, err := script.Run(ctx, s.client, []string{key},
		now,
		window,
		s.limit,
		requestID,
	).Result()

	if err != nil {
		return false, 0, 0, 0, err
	}

	val := res.([]any)
	allowed := val[0].(int64)
	totalCount := val[1].(int64)

	var oldestTimestamp int64

	if oldest, ok := val[2].([]any); ok && len(oldest) == 2 {
		scoreStr := oldest[1].(string)

		ts, err := strconv.ParseInt(scoreStr, 10, 64)
		if err != nil {
			return false, 0, 0, 0, err
		}

		oldestTimestamp = ts
	}

	remaining := max(int(s.limit-int(totalCount)), 0)
	reset := oldestTimestamp + window

	return allowed == 1, s.limit, remaining, reset, nil
}
