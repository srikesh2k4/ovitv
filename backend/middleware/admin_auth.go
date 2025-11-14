package middleware

import "net/http"

func AdminAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        session, _ := Store.Get(r, "admin-session")
        auth, ok := session.Values["auth"].(bool)

        if !ok || !auth {
            http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
            return
        }

        next.ServeHTTP(w, r)
    })
}
