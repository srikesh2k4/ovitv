package middleware

import (
    "net/http"
    "os"
)

func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        origin := r.Header.Get("Origin")

        // Allowed origins list
        allowedOrigins := map[string]bool{
            "http://localhost:3000":             true,
            "http://127.0.0.1:3000":             true,
            "http://localhost:5173":             true,
            "http://localhost:8080":             true,
            "https://ovitv.up.railway.app":      true, // CLIENT
            "https://ovitv-admin.up.railway.app": true, // ADMIN
            "https://ovitv.ddns.net":             true, // CUSTOM DOMAIN
        }

        // In production — only allow known domains
        if os.Getenv("ENV") == "production" {
            if allowedOrigins[origin] {
                w.Header().Set("Access-Control-Allow-Origin", origin)
            }
        } else {
            // In development — allow any origin
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }

        w.Header().Set("Vary", "Origin")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers",
            "Content-Type, Authorization, X-CSRF-Token, Accept, Origin, Upgrade, Connection")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
