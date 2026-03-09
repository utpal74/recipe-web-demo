package redis

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-demo/recipes-web/internal/platform/requestctx"
	"github.com/gin-gonic/gin"
)

func NewRateLimiter(limiter Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.ClientIP()

		if identity, ok := requestctx.GetAuthIdentity(ctx); ok {
			key = identity.UserID
		}

		key = fmt.Sprintf("rate_limit:%s:%v", limiter.Name(), key)

		allowed, limit, remaining, reset, err := limiter.Allow(ctx.Request.Context(), key)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		ctx.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		ctx.Header("X-RateLimit-Reset", strconv.FormatInt(reset, 10))

		if !allowed {
			retryAfter := max((reset-time.Now().UnixMilli())/1000, 0)
			ctx.Header("Retry-After", strconv.FormatInt(retryAfter, 10))
			
			ctx.JSON(http.StatusTooManyRequests, gin.H{"message": "too many requests"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
