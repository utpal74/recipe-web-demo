package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-demo/recipes-web/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

/*
TODO:
Small design note (not a correction)
Putting Claims and JWTOutput near AuthHandler is sensible for now, because:
they’re auth-related
they’re transport-facing
Later, you might split:
domain claims vs
HTTP response DTOs
But don’t worry about that yet.
*/

type Config struct {
	Secret string
	Issuer string
}

type AuthHandler struct {
	config Config
}

func New(config Config) *AuthHandler {
	return &AuthHandler{config: config}
}

func (ah *AuthHandler) SignInHandler(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "username or password is required",
		})
		return
	}

	if user.UserName != "admin" || user.Password != "password" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid username or password",
		})
		return
	}

	expiryAt := time.Now().Add(15 * time.Minute)
	token, err := ah.createToken(user.UserName, expiryAt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, JWTOutput{
		Token:   token,
		Expires: expiryAt,
	})
}

type Claims struct {
	UserName string `json:"userName"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expiresIn"`
}

func (ah *AuthHandler) createToken(userName string, expiry time.Time) (string, error) {
	claims := Claims{
		UserName: userName,
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			Issuer:    ah.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := ah.config.Secret
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET key not provided")
	}
	return token.SignedString([]byte(secretKey))
}
