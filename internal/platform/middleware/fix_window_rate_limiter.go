package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-demo/recipes-web/internal/platform/requestctx"
	"github.com/gin-gonic/gin"
)

type Counter struct {
	count       int
	windowStart time.Time
}

type FixWindowRateLimitingMiddleWare struct {
	totalLimit     int
	limiters       map[string]*Counter
	windowDuration time.Duration
	mu             sync.Mutex
}

func NewFixWindowRateLimitingMiddleWare(ctx context.Context, limit int, window time.Duration) *FixWindowRateLimitingMiddleWare {
	r := &FixWindowRateLimitingMiddleWare{
		totalLimit:     limit,
		limiters:       make(map[string]*Counter),
		windowDuration: window,
	}

	go r.cleanup(ctx, time.Minute*10)
	return r
}

func (f *FixWindowRateLimitingMiddleWare) cleanup(ctx context.Context, maxIdle time.Duration) {
	ticker := time.NewTicker(f.windowDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			f.mu.Lock()
			for k, counter := range f.limiters {
				if time.Since(counter.windowStart) >= maxIdle {
					delete(f.limiters, k)
				}
			}
			f.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (f *FixWindowRateLimitingMiddleWare) Allow(key string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	counter, ok := f.limiters[key]
	if !ok {
		counter = &Counter{
			count:       f.totalLimit,
			windowStart: time.Now(),
		}

		f.limiters[key] = counter
	}

	if time.Since(counter.windowStart) >= f.windowDuration {
		counter.count = f.totalLimit
		counter.windowStart = time.Now()
	}

	if counter.count > 0 {
		counter.count--
		return true
	}

	return false
}

func (f *FixWindowRateLimitingMiddleWare) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		key := ctx.ClientIP()

		if identity, ok := requestctx.GetAuthIdentity(ctx); ok {
			key = identity.UserID
		}

		if !f.Allow(key) {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
