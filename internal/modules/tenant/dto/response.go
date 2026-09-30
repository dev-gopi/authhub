package dto

import "github.com/google/uuid"

type CreateTenantResponse struct {
	ID uuid.UUID `json:"id"`

	APILabel string `json:"api_label"`

	DisplayName string `json:"display_name"`

	Status string `json:"status"`

	PrimaryAdmin PrimaryAdminResponse `json:"primary_admin"`

	DefaultUserPool UserPoolResponse `json:"default_user_pool"`
}

type PrimaryAdminResponse struct {
	ID uuid.UUID `json:"id"`

	Username string `json:"username"`

	Email string `json:"email"`
}

type UserPoolResponse struct {
	ID uuid.UUID `json:"id"`

	APILabel string `json:"api_label"`
}
