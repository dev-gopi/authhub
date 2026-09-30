package middleware

import "net/http"

// AuthMiddleware is intentionally fail-closed until the generic tenant-user
// authentication flow is implemented. Root authentication uses the dedicated
// rootauth middleware and does not depend on this placeholder.
func AuthMiddleware(_ http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "authentication middleware is not configured", http.StatusServiceUnavailable)
	})
}
