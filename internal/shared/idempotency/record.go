package idempotency

import "time"

// Status represents the state of action
type Status string

const (
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

// Record defines an idempotency record.
type Record struct {
	ID         string
	UserID     string
	Key        string
	Response   []byte
	Status     Status
	StatusCode int
	CreatedAt  time.Time
}
