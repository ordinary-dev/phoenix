package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.WithFields(log.Fields{
			"latency": time.Since(start),
			"method":  r.Method,
			"path":    r.URL.Path,
		}).Info("Request")
	})
}
