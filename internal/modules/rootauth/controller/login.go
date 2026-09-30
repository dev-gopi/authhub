package controller

import (
	"encoding/json"
	"net/http"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/dto"
	"github.com/dev-gopi/authhub/internal/shared/response"
)

func (c *Controller) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.LoginRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.Error(
			w,
			http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
		)
		return
	}

	if err :=
		c.validator.ValidateLogin(&req); err != nil {

		response.Error(
			w,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Invalid login request",
		)

		return
	}

	result, err := c.service.Login(
		r.Context(),
		req,
		requestMetadata(r),
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
