package service

import (
	"fmt"
	"time"

	"github.com/gin-demo/recipes-web/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

// TokenService - should responsible to create, validate, expire, re-new or refresh the token

type TokenService interface {
	CreateToken(identity domain.Identity, expiry time.Time) (string, error)
	ValidateToken(token string) (domain.Identity, error)
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

func (ts *jwtTokenService) CreateToken(identity domain.Identity, expiry time.Time) (string, error) {
	claims := jwtClaims{
		UserName: identity.UserName,
		Role:     identity.Role,
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

func (ts *jwtTokenService) ValidateToken(tokenString string) (domain.Identity, error) {
	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(ts.config.Secret), nil
	})

	if err != nil || !token.Valid {
		return domain.Identity{}, err
	}

	return domain.Identity{
		UserName: claims.UserName,
		Role:     claims.Role,
	}, nil
}
