package middleware

import (
	"fmt"
	"net/http"
)

func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Tenant middleware executed")
		next.ServeHTTP(w, r)
	})
}
