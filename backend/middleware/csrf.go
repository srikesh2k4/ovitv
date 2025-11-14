package middleware

import (
    "crypto/rand"
    "encoding/base64"
    "net/http"

    "github.com/gorilla/sessions"
)

var (
    // 64-byte auth key + 32-byte encryption key
    sessionAuthKey  = []byte("kjsdf83r9f8r3h9f8h39fh398h398h398h398h398h398h398h39")
    sessionEncryptKey = []byte("fj9e8h39fh398h398fh398fh398fh398")
    
    Store = sessions.NewCookieStore(sessionAuthKey, sessionEncryptKey)
)

func init() {
    Store.Options = &sessions.Options{
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // set TRUE in production
        SameSite: http.SameSiteLaxMode,
    }
}

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
