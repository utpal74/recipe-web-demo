package memory

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-demo/recipes-web/internal/platform/requestctx"
	"github.com/gin-gonic/gin"
)

type SlCounter struct {
	currentCount  int
	previousCount int
	windowStart   time.Time
	lastSeen      time.Time
}

type SlideWindowWeighedAvgMiddleWare struct {
	limit   int
	window  time.Duration
	counter map[string]*SlCounter
	mu      sync.Mutex
}

func NewSlideWindowWeighedAvgMiddleWare(ctx context.Context, limit int, window time.Duration) *SlideWindowWeighedAvgMiddleWare {
	sl := &SlideWindowWeighedAvgMiddleWare{
		limit:   limit,
		window:  window,
		counter: make(map[string]*SlCounter),
	}

	go sl.cleanup(ctx, 10*time.Minute)

	return sl
}

func (s *SlideWindowWeighedAvgMiddleWare) cleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			now := time.Now()

			for key, counter := range s.counter {
				if now.Sub(counter.lastSeen) > 2*s.window {
					delete(s.counter, key)
				}
			}

			s.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (s *SlideWindowWeighedAvgMiddleWare) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	c, ok := s.counter[key]
	if !ok {
		s.counter[key] = &SlCounter{
			currentCount:  1,
			previousCount: 0,
			windowStart:   now,
			lastSeen:      now,
		}
		return true
	}

	elapsed := now.Sub(c.windowStart)

	if elapsed >= s.window {
		windowPassed := int(elapsed / s.window)

		if windowPassed == 1 {
			c.previousCount = c.currentCount
		} else {
			c.previousCount = 0
		}

		c.currentCount = 0
		c.windowStart = now
		elapsed = 0
	}

	weight := 1 - float64(elapsed)/float64(s.window)
	effectiveCount := float64(c.currentCount) + float64(c.previousCount)*weight

	if effectiveCount >= float64(s.limit) {
		return false
	}

	c.currentCount++
	c.lastSeen = now

	return true
}

func (s *SlideWindowWeighedAvgMiddleWare) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.ClientIP()

		if identity, ok := requestctx.GetAuthIdentity(ctx.Request.Context()); ok {
			key = identity.UserID
		}

		if !s.Allow(key) {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
