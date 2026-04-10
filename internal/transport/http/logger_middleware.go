package http

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

			fields := logrus.Fields{
				"method":      r.Method,
				"uri":         r.RequestURI,
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
				"duration":    time.Since(start).String(),
				"status":      recorder.statusCode,
			}

			switch {
			case recorder.statusCode >= 500:
				logger.WithFields(fields).Error("HTTP request completed with error")
			case recorder.statusCode >= 400:
				logger.WithFields(fields).Warn("HTTP request completed with client error")
			default:
				logger.WithFields(fields).Info("HTTP request completed")
			}
		})
	}
}
