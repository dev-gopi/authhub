package constants

import "time"

const (
	SessionTokenType = "Bearer"

	SessionAbsoluteLifetime = 12 * time.Hour

	SessionIdleLifetime = 30 * time.Minute
)

const (
	FailureInvalidCredentials = "invalid_credentials"
	FailureDisabled           = "account_disabled"
	FailureRateLimited        = "rate_limited"
)
