package service

import "github.com/golang-jwt/jwt/v5"

type jwtClaims struct {
	UserName  string `json:"userName"`
	Role      string `json:"role"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}
