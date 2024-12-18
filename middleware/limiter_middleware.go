package middleware

import (
	"context"
	"net/http"

	"github.com/htm1000/rate-limiter/limiter"
)

func RateLimiterMiddleware(rateLimiter *limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()

			token := r.Header.Get("API_KEY")
			isToken := false
			key := r.RemoteAddr

			if token != "" {
				isToken = true
				key = token
			}

			allowed, err := rateLimiter.Allow(ctx, key, isToken)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "You have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
