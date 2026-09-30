package middleware

import "net/http"

// TenantMiddleware is intentionally fail-closed until tenant context
// resolution is implemented.
func TenantMiddleware(_ http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "tenant middleware is not configured", http.StatusServiceUnavailable)
	})
}
