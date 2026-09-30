package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/dev-gopi/authhub/internal/shared/response"

	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			defer func() {
				if recovered := recover(); recovered != nil {

					logger.Error(
						"panic recovered",
						zap.Any("panic", recovered),
						zap.ByteString("stack", debug.Stack()),
						zap.String(
							"request_id",
							GetRequestID(r.Context()),
						),
					)

					response.Error(
						w,
						http.StatusInternalServerError,
						"INTERNAL_ERROR",
						fmt.Sprintf(
							"internal server error",
						),
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
