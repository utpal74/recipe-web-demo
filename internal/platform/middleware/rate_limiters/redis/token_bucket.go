package redis

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

var tokenBucketScript = redis.NewScript(`
local key = KEYS[1]

local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4]
)
local data = redis.call('HMGET', key, "token", "timestamp")

local token = tonumber(data[1])
local last_refill = tonumber(data[2])

if token == nil then
   token = capacity
   last_refill = now
end

local delta = math.max(0, now-last_refill)
local refill = (delta/1000) * refill_rate

token = math.min(capacity, token + refill)
local allowed = token >= requested

if allowed then
   token = token - requested
end

redis.call("HSET", key, 
   "token", token, 
   "timestamp", now
)

redis.call("PEXPIRE", key, math.ceil((capacity/refill_rate)*1000))

return {allowed and 1 or 0, token}
`)

type tokenBucketLimiter struct {
	client   *redis.Client
	capacity int
	rate     int
}

// Allow implements [Limiter].
func (s *tokenBucketLimiter) Allow(ctx context.Context, key string) (bool, int, int, int64, error) {
	
	now := time.Now()

	res, err := tokenBucketScript.Run(ctx, s.client, 
		[]string{key}, 
		s.capacity, 
		s.rate, 
		now.UnixMilli(), 1).Result()
	if err != nil {
		return false, 0, 0, 0, err
	}

	val := res.([]any)
	allowed := val[0].(int64) == 1

	var token float64
	switch v:= val[1].(type) {
	case int64:
		token = float64(v)
	case float64:
		token = v
	default:
		return false, 0, 0, 0, fmt.Errorf("unexpected token type %T", v)
	}
	
	remaining := int(math.Floor(token))
	missing := 1 - token
	if missing < 0 {
		missing = 0
	}

	seconds := missing / float64(s.rate)

	reset := now.Add(time.Duration(seconds * float64(time.Second))).UnixMilli()
	return allowed, s.capacity, remaining, reset, nil
}

// Name implements [Limiter].
func (s *tokenBucketLimiter) Name() string {
	return "token_bucket"
}

func NewTokenBucketLimiter(client *redis.Client, capacity int, rate int) Limiter {
	return &tokenBucketLimiter{
		client:   client,
		capacity: capacity,
		rate:     rate,
	}
}
