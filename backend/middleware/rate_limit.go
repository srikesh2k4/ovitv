package middleware

import (
    "net/http"
    "sync"
    "time"

    "golang.org/x/time/rate"
    "ome-tv-pro/backend/utils"
)

var mu sync.Mutex
var visitors = make(map[string]*rate.Limiter)

func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        ip := utils.GetIP(r)

        mu.Lock()
        limiter, exists := visitors[ip]
        if !exists {
            limiter = rate.NewLimiter(rate.Every(time.Minute), 10)
            visitors[ip] = limiter
        }
        mu.Unlock()

        if !limiter.Allow() {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}
