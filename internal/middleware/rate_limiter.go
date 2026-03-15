package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter stores rate limit information
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Clean up old entries periodically
	go rl.cleanup()

	return rl
}

// cleanup removes old entries from the rate limiter
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, times := range rl.requests {
			// Remove old timestamps
			validTimes := []time.Time{}
			for _, t := range times {
				if now.Sub(t) < rl.window {
					validTimes = append(validTimes, t)
				}
			}

			if len(validTimes) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validTimes
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Get existing requests for this key
	times, exists := rl.requests[key]
	if !exists {
		rl.requests[key] = []time.Time{now}
		return true
	}

	// Remove old timestamps
	validTimes := []time.Time{}
	for _, t := range times {
		if now.Sub(t) < rl.window {
			validTimes = append(validTimes, t)
		}
	}

	// Check if we're at the limit
	if len(validTimes) >= rl.limit {
		rl.requests[key] = validTimes
		return false
	}

	// Add new request
	validTimes = append(validTimes, now)
	rl.requests[key] = validTimes

	return true
}

// GetRemaining returns the number of remaining requests
func (rl *RateLimiter) GetRemaining(key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	times, exists := rl.requests[key]
	if !exists {
		return rl.limit
	}

	now := time.Now()
	validCount := 0
	for _, t := range times {
		if now.Sub(t) < rl.window {
			validCount++
		}
	}

	return rl.limit - validCount
}

var (
	// AuthRateLimiter limits authentication endpoints (5 requests per minute per IP)
	AuthRateLimiter = NewRateLimiter(5, time.Minute)

	// GeneralRateLimiter limits general endpoints (100 requests per minute per IP)
	GeneralRateLimiter = NewRateLimiter(100, time.Minute)
)

// RateLimitMiddleware provides rate limiting for endpoints
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP address)
		clientIP := c.ClientIP()

		// Check if request should be allowed
		if !limiter.Allow(clientIP) {
			remaining := limiter.GetRemaining(clientIP)
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.limit))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Header("X-RateLimit-Reset", time.Now().Add(limiter.window).Format(time.RFC3339))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":      "Too many requests. Please try again later.",
				"request_id": c.GetString("request_id"),
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		remaining := limiter.GetRemaining(clientIP)
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		c.Next()
	}
}
