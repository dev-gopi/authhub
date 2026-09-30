package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

const RequestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		requestID := r.Header.Get(RequestIDHeader)

		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx := context.WithValue(
			r.Context(),
			RequestIDKey,
			requestID,
		)

		w.Header().Set(RequestIDHeader, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	value, _ := ctx.Value(RequestIDKey).(string)

	return value
}
