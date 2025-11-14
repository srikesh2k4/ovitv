package utils

import (
    "net"
    "net/http"
    "strings"
)

func GetIP(r *http.Request) string {
    // 1. Check X-Forwarded-For (can contain multiple IPs)
    fwd := r.Header.Get("X-Forwarded-For")
    if fwd != "" {
        // Extract the first valid IP
        parts := strings.Split(fwd, ",")
        ip := strings.TrimSpace(parts[0])
        if net.ParseIP(ip) != nil {
            return ip
        }
    }

    // 2. Check X-Real-IP
    realIP := r.Header.Get("X-Real-IP")
    if realIP != "" && net.ParseIP(realIP) != nil {
        return realIP
    }

    // 3. Fallback: use RemoteAddr
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err == nil && net.ParseIP(ip) != nil {
        return ip
    }

    // 4. Last fallback (uncommon)
    return r.RemoteAddr
}
