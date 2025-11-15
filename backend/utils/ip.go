package utils

import (
    "net"
    "net/http"
    "strings"
)

// GetIP extracts the real client IP from headers or fallback connection.
func GetIP(r *http.Request) string {
    // 1. Check X-Forwarded-For (can contain multiple IPs: "client, proxy1, proxy2")
    if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
        parts := strings.Split(fwd, ",")
        ip := strings.TrimSpace(parts[0])
        if net.ParseIP(ip) != nil {
            return ip
        }
    }

    // 2. Check X-Real-IP
    if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
        if net.ParseIP(realIP) != nil {
            return realIP
        }
    }

    // 3. Fallback to RemoteAddr
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err == nil && net.ParseIP(ip) != nil {
        return ip
    }

    // 4. Very last fallback (rare cases)
    return r.RemoteAddr
}
