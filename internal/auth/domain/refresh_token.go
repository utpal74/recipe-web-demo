package domain

import (
	"time"
)

type TokenID string

type UserID string

type RefreshToken struct {
	ID        TokenID
	UserID    UserID
	TokenHash string
	Revoked   bool
	CreatedAt time.Time
	ExpiresAt time.Time

	// TODO - provid support for this, for showing active sessions and revoking active devices.
	// UserAgent string
	// IP string
}
