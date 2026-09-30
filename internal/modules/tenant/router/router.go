package router

import (
	tenantcontroller "github.com/dev-gopi/authhub/internal/modules/tenant/controller"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Register(r chi.Router, controller *tenantcontroller.Controller, authMiddleware func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/api/v1/root/tenants", controller.Create)
	})
}
