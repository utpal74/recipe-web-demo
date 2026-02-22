package middleware

import (
	"net/http"

	"github.com/gin-demo/recipes-web/internal/platform/requestctx"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitingMiddleWare defines a rate limiter.
type RateLimitingMiddleWare struct {
	limit    rate.Limit
	burst    int
	limiters map[string]*rate.Limiter
}

// NewRateLimitMiddleware creates a rate limiter.
func NewRateLimitMiddleware(limit rate.Limit, burst int) *RateLimitingMiddleWare {
	return &RateLimitingMiddleWare{
		limit:    limit,
		burst:    burst,
		limiters: make(map[string]*rate.Limiter),
	}
}

func (r *RateLimitingMiddleWare) getLimiter(key string) *rate.Limiter {
	l, ok := r.limiters[key]
	if !ok {
		l = rate.NewLimiter(r.limit, r.burst)
		r.limiters[key] = l
	}

	return l
}

func (r *RateLimitingMiddleWare) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		key := ctx.ClientIP()

		if identity, ok := requestctx.GetAuthIdentity(ctx); ok {
			key = identity.UserID
		}

		limiter := r.getLimiter(key)

		if !limiter.Allow() {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
