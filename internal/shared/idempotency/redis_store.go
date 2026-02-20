package idempotency

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisStore struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisStore - creates a new redis store
func NewRedisStore(client *redis.Client, ttl time.Duration) IdempotencyStore {
	return &redisStore{
		client: client,
		ttl: ttl,
	}
}

func idempotencyKey(userID, key string) string {
	return "Idempotency:" + userID + key
}

// Get implements [IdempotencyStore].
func (r *redisStore) Get(ctx context.Context, userID string, key string) (*Record, error) {
	bytes, err := r.client.Get(ctx, idempotencyKey(userID, key)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}

	var record Record
	if err := json.Unmarshal(bytes, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

// Set implements [IdempotencyStore].
func (r *redisStore) Set(ctx context.Context, record *Record) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return r.client.Set(ctx,
		idempotencyKey(record.UserID, record.Key),
		data,
		r.ttl,
	).Err()
}

// TryLock implements [IdempotencyStore].
func (r *redisStore) TryLock(ctx context.Context, record *Record) (bool, error) {
	record.Status = StatusProcessing

	data, err := json.Marshal(record)
	if err != nil {
		return false, err
	}

	return r.client.SetNX(ctx,
		idempotencyKey(record.UserID, record.Key),
		data,
		r.ttl,
	).Result()
}
