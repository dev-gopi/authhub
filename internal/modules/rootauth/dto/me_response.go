package dto

import "github.com/google/uuid"

type MeResponse struct {
	ID uuid.UUID `json:"id"`

	Username string `json:"username"`

	Email string `json:"email"`

	DisplayName string `json:"display_name"`

	IsRootAdmin bool `json:"is_root_admin"`
}
