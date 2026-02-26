package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-demo/recipes-web/internal/platform/requestctx"
	"github.com/gin-gonic/gin"
)

type RequestLog struct {
	timeStamp time.Time
}

type SlideWindowLogMiddleWare struct {
	window time.Duration
	limit  int
	logs   map[string][]RequestLog
	mu     sync.Mutex
}

func NewSlideWindowLogMiddleWare(
	ctx context.Context,
	limit int,
	window time.Duration,
) *SlideWindowLogMiddleWare {
	m := &SlideWindowLogMiddleWare{
		window: window,
		limit:  limit,
		logs:   make(map[string][]RequestLog),
	}

	go m.cleanup(ctx, 1*time.Minute)
	return m
}

func (s *SlideWindowLogMiddleWare) cleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			s.mu.Lock()

			for key, log := range s.logs {

				valid := log[:0]

				for _, l := range log {
					if now.Sub(l.timeStamp) < s.window {
						valid = append(valid, l)
					}
				}

				if len(valid) == 0 {
					delete(s.logs, key)
				} else {
					s.logs[key] = valid
				}
			}
			s.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (s *SlideWindowLogMiddleWare) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	logs, ok := s.logs[key]
	if !ok {
		reqLog := RequestLog{timeStamp: now}
		s.logs[key] = []RequestLog{reqLog}
		return true
	}

	valid := logs[:0]
	for _, log := range logs {

		if time.Since(log.timeStamp) < s.window {
			valid = append(valid, log)
		}
	}

	if len(valid) >= s.limit {
		return false
	}

	valid = append(valid, RequestLog{timeStamp: now})
	s.logs[key] = valid
	return true
}

func (s *SlideWindowLogMiddleWare) Handle() gin.HandlerFunc {
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
