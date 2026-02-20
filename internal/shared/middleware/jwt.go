package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-demo/recipes-web/internal/auth/service"
	"github.com/gin-gonic/gin"
)

// JWTMiddleware handles JWT authentication for HTTP requests.
type JWTMiddleware struct {
	tokenService service.TokenService
}

// NewAuthMiddleWare creates a new JWTMiddleware for authentication.
func NewAuthMiddleWare(tokenService service.TokenService) *JWTMiddleware {
	return &JWTMiddleware{
		tokenService: tokenService,
	}
}

func (au *JWTMiddleware) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		identity, err := au.tokenService.ValidateAccessToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		newCtx := context.WithValue(ctx.Request.Context(), "authIdentity", identity)
		ctx.Request = ctx.Request.WithContext(newCtx)
		// ctx.Set("authIdentity", identity)

		idemKey := ctx.GetHeader("Idempotency-Key")
		if idemKey != "" {
			// ctx.Set("idempotencyKey", idemKey)
			newCtx := context.WithValue(ctx.Request.Context(), "idempotencyKey", idemKey)
			ctx.Request = ctx.Request.WithContext(newCtx)
		}

		ctx.Next()
	}
}
