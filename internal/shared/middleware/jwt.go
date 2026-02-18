package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-demo/recipes-web/internal/auth/service"
	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	tokenService service.TokenService
}

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

		ctx.Set("identity", identity)
		ctx.Next()
	}
}
