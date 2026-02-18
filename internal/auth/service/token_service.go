package service

import (
	"fmt"
	"time"

	"github.com/gin-demo/recipes-web/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

// TokenService - should responsible to create, validate, expire, re-new or refresh the token

type TokenService interface {
	CreateAccessToken(identity domain.Identity, expiry time.Time) (string, error)
	ValidateAccessToken(token string) (domain.Identity, error)

	CreateRefreshToken(identity domain.Identity, expiry time.Time) (string, error)
	ValidateRefreshToken(token string) (domain.Identity, error)
}

type jwtTokenService struct {
	config Config
}

type Config struct {
	Secret string
	Issuer string
}

func NewJwtTokenService(config Config) TokenService {
	return &jwtTokenService{config: config}
}

func (ts *jwtTokenService) CreateAccessToken(identity domain.Identity, expiry time.Time) (string, error) {
	claims := jwtClaims{
		UserID:    identity.UserID,
		Role:      identity.Role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			Issuer:    ts.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := ts.config.Secret
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET key not provided")
	}
	return token.SignedString([]byte(secretKey))
}

func (ts *jwtTokenService) ValidateAccessToken(tokenString string) (domain.Identity, error) {
	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(ts.config.Secret), nil
	})

	if err != nil || !token.Valid {
		return domain.Identity{}, err
	}

	return domain.Identity{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}

// CreateRefreshToken implements [TokenService].
func (ts *jwtTokenService) CreateRefreshToken(identity domain.Identity, expiry time.Time) (string, error) {
	claims := jwtClaims{
		UserID:    identity.UserID,
		Role:      identity.Role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			Issuer:    ts.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := ts.config.Secret
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET key not provided")
	}
	return token.SignedString([]byte(secretKey))
}

// ValidateRefreshToken implements [TokenService].
func (ts *jwtTokenService) ValidateRefreshToken(tokenString string) (domain.Identity, error) {
	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(ts.config.Secret), nil
	})

	if err != nil || !token.Valid {
		return domain.Identity{}, err
	}

	if claims.TokenType != "refresh" {
		return domain.Identity{}, fmt.Errorf("invlid token type: %w", err)
	}

	return domain.Identity{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
