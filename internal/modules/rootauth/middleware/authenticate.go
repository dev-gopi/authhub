package middleware

import (
	"net/http"
	"strings"

	rootservice "github.com/dev-gopi/authhub/internal/modules/rootauth/service"
	"github.com/dev-gopi/authhub/internal/shared/response"
)

type Middleware struct {
	service rootservice.Interface
}

func New(
	service rootservice.Interface,
) *Middleware {
	return &Middleware{
		service: service,
	}
}

func (m *Middleware) Authenticate(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		header := strings.TrimSpace(
			r.Header.Get("Authorization"),
		)

		if header == "" {
			response.Error(
				w,
				http.StatusUnauthorized,
				"UNAUTHORIZED",
				"Authentication required",
			)
			return
		}

		parts := strings.Fields(header)

		if len(parts) != 2 ||
			!strings.EqualFold(
				parts[0],
				"Bearer",
			) {
			response.Error(
				w,
				http.StatusUnauthorized,
				"UNAUTHORIZED",
				"Invalid authorization header",
			)
			return
		}

		auth, err :=
			m.service.Authenticate(
				r.Context(),
				parts[1],
			)

		if err != nil {
			response.Error(
				w,
				http.StatusUnauthorized,
				"UNAUTHORIZED",
				"Invalid or expired session",
			)
			return
		}

		ctx := WithRootAuthContext(
			r.Context(),
			auth,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}
