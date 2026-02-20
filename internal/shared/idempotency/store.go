package idempotency

import "context"

// IdempotencyStore defines the contract for idempotency.
type IdempotencyStore interface {
	Get(ctx context.Context, userID, key string) (*Record, error)
	Set(ctx context.Context, record *Record) error
	TryLock(ctx context.Context, record *Record) (bool, error)
}
