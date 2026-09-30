package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func RequestLogger(
	logger *zap.Logger,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			start := time.Now()

			writer := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(writer, r)

			logger.Info(
				"http request",

				zap.String(
					"method",
					r.Method,
				),

				zap.String(
					"path",
					r.URL.Path,
				),

				zap.Int(
					"status",
					writer.status,
				),

				zap.Duration(
					"duration",
					time.Since(start),
				),

				zap.String(
					"request_id",
					GetRequestID(r.Context()),
				),

				zap.String(
					"correlation_id",
					GetCorrelationID(r.Context()),
				),
			)
		})
	}
}
