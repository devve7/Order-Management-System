// Package middleware ...
package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleWare(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			recorder := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(recorder, r)

			logger.WithFields(logrus.Fields{
				"method":      r.Method,
				"uri":         r.RequestURI,
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
				"duration":    time.Since(start).String(),
				"status":      recorder.statusCode,
			}).Info("HTTP request completed")
		})
	}
}
