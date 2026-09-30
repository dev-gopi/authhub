package controller

import (
	"errors"
	"net"
	"net/http"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/entity"
	rootservice "github.com/dev-gopi/authhub/internal/modules/rootauth/service"
	rootvalidator "github.com/dev-gopi/authhub/internal/modules/rootauth/validator"
	"github.com/dev-gopi/authhub/internal/shared/response"
)

type Controller struct {
	service rootservice.Interface

	validator *rootvalidator.Validator
}

func New(
	service rootservice.Interface,
	validator *rootvalidator.Validator,
) *Controller {
	return &Controller{
		service:   service,
		validator: validator,
	}
}

func requestMetadata(
	r *http.Request,
) entity.RequestMetadata {
	host := r.RemoteAddr

	if parsedHost, _, err :=
		net.SplitHostPort(r.RemoteAddr); err == nil {
		host = parsedHost
	}

	return entity.RequestMetadata{
		IPAddress: host,

		UserAgent: r.UserAgent(),
	}
}

func handleError(
	w http.ResponseWriter,
	err error,
) {
	switch {
	case errors.Is(
		err,
		entity.ErrInvalidCredentials,
	):
		response.Error(
			w,
			http.StatusUnauthorized,
			"INVALID_CREDENTIALS",
			"Invalid credentials",
		)

	case errors.Is(
		err,
		entity.ErrUnauthorized,
	),
		errors.Is(
			err,
			entity.ErrSessionExpired,
		):
		response.Error(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Authentication required",
		)

	case errors.Is(
		err,
		entity.ErrForbidden,
	):
		response.Error(
			w,
			http.StatusForbidden,
			"FORBIDDEN",
			"Access denied",
		)

	case errors.Is(
		err,
		entity.ErrRateLimited,
	):
		response.Error(
			w,
			http.StatusTooManyRequests,
			"RATE_LIMITED",
			"Too many login attempts",
		)

	default:
		response.Error(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Internal server error",
		)
	}
}
