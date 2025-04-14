package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func (m *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		slog.Info(
			"request",
			"latency", time.Since(start),
			"method", r.Method,
			"path", r.URL.Path,
		)
	})
}
