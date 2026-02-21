package service

import "time"

// AuthConfig - represents the TTL configuration.
type AuthConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}
