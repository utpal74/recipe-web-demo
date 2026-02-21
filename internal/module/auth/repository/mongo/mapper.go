package mongo

import "github.com/gin-demo/recipes-web/internal/module/auth/domain"

func toDomain(doc refreshTokenDoc) *domain.RefreshToken {
	return &domain.RefreshToken{
		ID:        domain.TokenID(doc.ID),
		UserID:    domain.UserID(doc.UserID),
		TokenHash: doc.TokenHash,
		Revoked:   doc.Revoked,
		CreatedAt: doc.CreatedAt,
		ExpiresAt: doc.ExpiresAt,
	}
}

func fromDomain(token *domain.RefreshToken) refreshTokenDoc {
	return refreshTokenDoc{
		ID:        string(token.ID),
		UserID:    string(token.UserID),
		TokenHash: token.TokenHash,
		Revoked:   token.Revoked,
		CreatedAt: token.CreatedAt,
		ExpiresAt: token.ExpiresAt,
	}
}
