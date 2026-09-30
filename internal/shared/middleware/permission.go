package middleware

import "net/http"

// PermissionMiddleware is intentionally fail-closed until tenant permission
// evaluation is implemented. It must never behave as an allow-all placeholder.
func PermissionMiddleware(_ http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "permission middleware is not configured", http.StatusServiceUnavailable)
	})
}
