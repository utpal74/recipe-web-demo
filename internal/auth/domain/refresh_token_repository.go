package domain

import "context"

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *RefreshToken) error
	FindByTokenHash(ctx context.Context, hash string) (*RefreshToken, error)
	DeleteByID(ctx context.Context, id TokenID) error
}
