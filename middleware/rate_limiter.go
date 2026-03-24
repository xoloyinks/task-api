// middleware/rate_limiter.go
package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// holds a limiter per IP address
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

func getClientLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	// if client exists return their limiter
	if c, exists := clients[ip]; exists {
		c.lastSeen = time.Now()
		return c.limiter
	}

	// create a new limiter for new clients
	// rate.NewLimiter(r, b)
	// r = requests per second
	// b = burst size (max requests at once)
	limiter := rate.NewLimiter(2, 5) // 2 requests/sec, burst of 5
	clients[ip] = &client{
		limiter:  limiter,
		lastSeen: time.Now(),
	}

	return limiter
}

// clean up clients that haven't been seen in 3 minutes
// runs in background to prevent memory leak
func cleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func RateLimiterMiddleware(next http.Handler) http.Handler {
	// start cleanup goroutine
	go cleanupClients()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		limiter := getClientLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, `{"error": "too many requests"}`, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
