package controller

import (
	"net/http"

	rootmiddleware "github.com/dev-gopi/authhub/internal/modules/rootauth/middleware"
	"github.com/dev-gopi/authhub/internal/shared/response"
)

func (c *Controller) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {
	auth, ok :=
		rootmiddleware.FromRootAuthContext(
			r.Context(),
		)

	if !ok {
		response.Error(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Authentication required",
		)
		return
	}

	if err := c.service.Logout(
		r.Context(),
		auth,
	); err != nil {
		handleError(w, err)
		return
	}

	response.Success(
		w,
		http.StatusOK,
		map[string]string{
			"message": "Logged out successfully",
		},
	)
}

func (c *Controller) LogoutAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	auth, ok :=
		rootmiddleware.FromRootAuthContext(
			r.Context(),
		)

	if !ok {
		response.Error(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Authentication required",
		)
		return
	}

	if err := c.service.LogoutAll(
		r.Context(),
		auth,
	); err != nil {
		handleError(w, err)
		return
	}

	response.Success(
		w,
		http.StatusOK,
		map[string]string{
			"message": "All sessions logged out successfully",
		},
	)
}
