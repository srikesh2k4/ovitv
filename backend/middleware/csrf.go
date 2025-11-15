package middleware

import (
    "crypto/rand"
    "encoding/base64"
    "net/http"
    "os"

    "github.com/gorilla/sessions"
)

// -------- Session Keys (LOAD FROM ENV) --------

func getEnvKey(name string, size int) []byte {
    val := os.Getenv(name)
    if len(val) < size {
        // auto generate fallback key (only for dev)
        b := make([]byte, size)
        rand.Read(b)
        return b
    }
    return []byte(val)
}

var (
    sessionAuthKey    = getEnvKey("SESSION_AUTH_KEY", 32)    // 32 bytes
    sessionEncryptKey = getEnvKey("SESSION_ENCRYPT_KEY", 32) // 32 bytes

    Store = sessions.NewCookieStore(sessionAuthKey, sessionEncryptKey)
)

func init() {
    Store.Options = &sessions.Options{
        Path:     "/",
        HttpOnly: true,
        Secure:   true,                          // REQUIRED for Railway (HTTPS)
        SameSite: http.SameSiteNoneMode,         // works with cross-site admin panel
    }
}

// -------- CSRF Helpers --------

func GenerateCSRFToken() (string, error) {
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(b), nil
}

func SetCSRF(w http.ResponseWriter, r *http.Request) (string, error) {
    session, _ := Store.Get(r, "admin-session")

    token, err := GenerateCSRFToken()
    if err != nil {
        return "", err
    }

    session.Values["csrf"] = token
    session.Save(r, w)

    return token, nil
}

func VerifyCSRF(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            next.ServeHTTP(w, r)
            return
        }

        session, _ := Store.Get(r, "admin-session")
        stored := session.Values["csrf"]

        sent := r.Header.Get("X-CSRF-Token")
        if sent == "" {
            sent = r.FormValue("csrf")
        }

        if stored == nil || sent != stored.(string) {
            http.Error(w, "Invalid CSRF token", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}
