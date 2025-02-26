package middleware

import (
	"net/http"
	"golang.org/x/time/rate"
)

func RateLimit(limit float64, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(limit), burst) // e.g., 10 req/s, burst of 20
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}