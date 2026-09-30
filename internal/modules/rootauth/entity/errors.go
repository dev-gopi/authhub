package entity

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrUnauthorized = errors.New("unauthorized")

	ErrForbidden = errors.New("forbidden")

	ErrSessionExpired = errors.New("session expired")

	ErrRateLimited = errors.New("too many login attempts")

	ErrAccountUnavailable = errors.New("account unavailable")
)
