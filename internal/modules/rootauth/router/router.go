package router

import (
	"net/http"

	rootcontroller "github.com/dev-gopi/authhub/internal/modules/rootauth/controller"

	"github.com/go-chi/chi/v5"
)

func Register(
	r chi.Router,
	controller *rootcontroller.Controller,
	authMiddleware func(http.Handler) http.Handler,
) {
	r.Route(
		"/api/v1/root/auth",
		func(r chi.Router) {
			r.Post(
				"/login",
				controller.Login,
			)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware)

				r.Get(
					"/me",
					controller.Me,
				)

				r.Post(
					"/logout",
					controller.Logout,
				)

				r.Post(
					"/logout-all",
					controller.LogoutAll,
				)
			})
		},
	)
}
