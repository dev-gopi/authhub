package dto

import "time"

type LoginResponse struct {
	SessionToken string `json:"session_token"`

	TokenType string `json:"token_type"`

	ExpiresAt time.Time `json:"expires_at"`
}
