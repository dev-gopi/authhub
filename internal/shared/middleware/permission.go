package middleware

import (
	"fmt"
	"net/http"
)

func PermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Permission middleware executed")
		next.ServeHTTP(w, r)
	})
}
