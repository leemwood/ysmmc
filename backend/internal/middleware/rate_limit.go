package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ysmmc/backend/pkg/response"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
}

type visitor struct {
	lastSeen time.Time
	count    int
}

type RateLimitConfig struct {
	Requests   int
	Window     time.Duration
	KeyFunc    func(*gin.Context) string
}

var limiter = &RateLimiter{
	visitors: make(map[string]*visitor),
}

func init() {
	go cleanupVisitors()
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		limiter.mu.Lock()
		for ip, v := range limiter.visitors {
			if time.Since(v.lastSeen) > time.Hour {
				delete(limiter.visitors, ip)
			}
		}
		limiter.mu.Unlock()
	}
}

func getVisitor(key string) *visitor {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	v, exists := limiter.visitors[key]
	if !exists {
		v = &visitor{
			lastSeen: time.Now(),
			count:    0,
		}
		limiter.visitors[key] = v
	}
	return v
}

func RateLimit(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := config.KeyFunc(c)
		if key == "" {
			c.Next()
			return
		}

		v := getVisitor(key)
		
		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		now := time.Now()
		
		if now.Sub(v.lastSeen) > config.Window {
			v.count = 0
			v.lastSeen = now
		}

		v.count++
		v.lastSeen = now

		if v.count > config.Requests {
			response.TooManyRequests(c, "too many requests, please try again later")
			c.Abort()
			return
		}

		c.Next()
	}
}

func IPKeyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func LoginRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Requests: 5,
		Window:   time.Minute,
		KeyFunc:  IPKeyFunc,
	})
}

func RegisterRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Requests: 3,
		Window:   time.Hour,
		KeyFunc:  IPKeyFunc,
	})
}

func ForgotPasswordRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Requests: 3,
		Window:   time.Hour,
		KeyFunc:  IPKeyFunc,
	})
}

func GlobalRateLimit() gin.HandlerFunc {
	return RateLimit(RateLimitConfig{
		Requests: 100,
		Window:   time.Second,
		KeyFunc:  IPKeyFunc,
	})
}
