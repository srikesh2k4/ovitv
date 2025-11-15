package middleware

import (
    "net/http"
    "os"
)

func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

        // Enable HSTS only in production
        if os.Getenv("ENV") == "production" {
            w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }

        next.ServeHTTP(w, r)
    })
}

func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        origin := r.Header.Get("Origin")

        // Development allowed origins
        if origin == "http://localhost:3000" ||
            origin == "http://127.0.0.1:3000" ||
            origin == "http://localhost:5173" ||
            origin == "http://localhost:8080" ||
            os.Getenv("ENV") != "production" {

            w.Header().Set("Access-Control-Allow-Origin", origin)

        } else {
            // Production domain
            w.Header().Set("Access-Control-Allow-Origin", "https://ovitv.ddns.net")
        }

        w.Header().Set("Vary", "Origin")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers",
            "Content-Type, Authorization, X-CSRF-Token, Accept, Origin, Upgrade, Connection")

        // Preflight
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
