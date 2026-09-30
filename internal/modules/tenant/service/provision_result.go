package service

import "github.com/google/uuid"

type provisionResult struct {
	UserPoolID         uuid.UUID
	PrimaryAdminRoleID uuid.UUID
	PrimaryAdminID     uuid.UUID
}
