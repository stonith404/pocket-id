package middleware

import (
	"golang-rest-api-template/internal/common"
	"golang-rest-api-template/internal/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter is a Gin middleware for rate limiting based on client IP
func RateLimiter(limit rate.Limit, burst int) gin.HandlerFunc {
	// Start the cleanup routine
	go cleanupClients()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Skip rate limiting for localhost and test environment
		// If the client ip is localhost the request comes from the frontend
		if ip == "127.0.0.1" || ip == "::1" || common.EnvConfig.AppEnv == "test" {
			c.Next()
			return
		}

		limiter := getLimiter(ip, limit, burst)
		if !limiter.Allow() {
			utils.HandlerError(c, http.StatusTooManyRequests, "Too many requests. Please wait a while before trying again.")
			c.Abort()
			return
		}

		c.Next()
	}
}

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Map to store the rate limiters per IP
var clients = make(map[string]*client)
var mu sync.Mutex

// Cleanup routine to remove stale clients that haven't been seen for a while
func cleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// getLimiter retrieves the rate limiter for a given IP address, creating one if it doesn't exist
func getLimiter(ip string, limit rate.Limit, burst int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if client, exists := clients[ip]; exists {
		client.lastSeen = time.Now()
		return client.limiter
	}

	limiter := rate.NewLimiter(limit, burst)
	clients[ip] = &client{limiter: limiter, lastSeen: time.Now()}
	return limiter
}
