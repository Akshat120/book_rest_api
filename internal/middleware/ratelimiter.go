package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

var rateLimiter *RateLimiter

type RateLimiter struct {
	Enabled bool
	Rate    float64
	Burst   int
	rl      *rate.Limiter
}

func InitRateLimiter(rps float64, burst int) *RateLimiter {
	if rateLimiter != nil {
		return rateLimiter
	}

	rateLimiter = &RateLimiter{
		Enabled: true,
		Rate:    float64(rps),
		Burst:   burst,
		rl:      rate.NewLimiter(rate.Limit(rps), burst),
	}

	return rateLimiter
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rateLimiter.Enabled && !rateLimiter.rl.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
