package memory

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-demo/recipes-web/internal/platform/requestctx"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type clientRateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimitingMiddleWare defines a rate limiter.
type RateLimitingMiddleWare struct {
	limit   rate.Limit
	burst   int
	clients map[string]*clientRateLimiter
	mu      sync.RWMutex
}

// NewRateLimitMiddleware creates a rate limiter.
func NewRateLimitMiddleware(ctx context.Context, limit rate.Limit, burst int) *RateLimitingMiddleWare {
	r := &RateLimitingMiddleWare{
		limit:   limit,
		burst:   burst,
		clients: make(map[string]*clientRateLimiter),
	}

	go r.cleanup(ctx, 5*time.Minute, 10*time.Minute)
	return r
}

func (r *RateLimitingMiddleWare) getLimiter(key string) *rate.Limiter {
	r.mu.RLock()
	cl, ok := r.clients[key]
	r.mu.RUnlock()

	if ok {
		r.mu.Lock()
		cl.lastSeen = time.Now()
		r.mu.Unlock()
		return cl.limiter
	}

	// if doesn't exist - get full lock
	r.mu.Lock()
	defer r.mu.Unlock()

	cl, ok = r.clients[key]
	if ok {
		cl.lastSeen = time.Now()
		return cl.limiter
	}

	limiter := rate.NewLimiter(r.limit, r.burst)
	r.clients[key] = &clientRateLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter
}

func (r *RateLimitingMiddleWare) cleanup(ctx context.Context, interval, maxIdle time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.mu.Lock()
			for k, cl := range r.clients {
				if time.Since(cl.lastSeen) > maxIdle {
					delete(r.clients, k)
				}
			}
			r.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
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
