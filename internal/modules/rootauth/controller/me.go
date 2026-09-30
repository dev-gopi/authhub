package controller

import (
	"net/http"

	rootmiddleware "github.com/dev-gopi/authhub/internal/modules/rootauth/middleware"
	"github.com/dev-gopi/authhub/internal/shared/response"
)

func (c *Controller) Me(
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

	result, err := c.service.Me(
		r.Context(),
		auth,
	)

	if err != nil {
		handleError(w, err)
		return
	}

	response.Success(
		w,
		http.StatusOK,
		result,
	)
}
