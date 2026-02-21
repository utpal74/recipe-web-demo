package service

import "github.com/golang-jwt/jwt/v5"

type jwtClaims struct {
	UserID    string `json:"id"`
	Role      string `json:"role"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}
