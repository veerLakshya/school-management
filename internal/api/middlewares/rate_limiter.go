package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu        sync.Mutex
	visitors  map[string]int
	limit     int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
	// start the reset routine
	go rl.ResetVisitorCount()
	return rl
}

func (rl *rateLimiter) ResetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) Middleware(next http.Handler) http.Handler {
	fmt.Println("Rate Limiter Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Rate Limiter Middleware begin returned...")
		rl.mu.Lock()
		defer rl.mu.Unlock()

		// get requestIP and increment its count
		visitorIP := r.RemoteAddr // We might extract IP in a more sophisticated way
		rl.visitors[visitorIP]++
		log.Printf("RATE LIMITER middleware - Visitor count of %v is %v\n", visitorIP, rl.visitors[visitorIP])

		if rl.visitors[visitorIP] > rl.limit {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
		fmt.Println("Rate Limiter Middleware ends...")
	})
}

// mutexes prevents race conditions by making sure that “Only one goroutine may execute the code between these two lines at a time.”
