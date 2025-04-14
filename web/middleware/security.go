package middleware

import (
	"net/http"
)

// Adds several headers to the response to improve security.
// For example, headers prevent embedding a site and passing information about the referrer.
func (m *Middleware) SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Content-Security-Policy", "script-src 'self'; ")

		next.ServeHTTP(w, r)
	})
}
