package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, server Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		httpVersion := r.Header.Get("Proto")
		userAgent := r.Header.Get("User-Agent")
		duration := time.Since(start)

		logData := fmt.Sprintf(
			"%s %s %s %s %s %ds",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			httpVersion,
			userAgent,
			duration,
		)

		server.logger.Info(logData)
	})
}
