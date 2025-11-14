package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// Wraps multiple middleware functions
func Chain(mw ...Middleware) Middleware {
    return func(final http.Handler) http.Handler {
        for i := len(mw) - 1; i >= 0; i-- {
            final = mw[i](final)
        }
        return final
    }
}
