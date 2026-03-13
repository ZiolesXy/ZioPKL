package middleware

import (
	"net/http"
	"sync"
	"time"
	"voca-plane/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	Visitors = make(map[string]*Visitor)
	mu      sync.Mutex
)

type Visitor struct{
	limiter *rate.Limiter
	lastSeen time.Time
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		if _, exists := Visitors[ip]; !exists {
			Visitors[ip] = &Visitor{limiter: rate.NewLimiter(30, 20)}
		}
		v := Visitors[ip]
		mu.Unlock()

		if !v.limiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, "too many request")
			c.Abort()
			return
		}

		c.Next()
	}
}