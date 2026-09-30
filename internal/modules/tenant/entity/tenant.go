package entity

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID uuid.UUID

	APILabel string

	DisplayName string

	Status Status

	DefaultLocale string

	SuspendedAt *time.Time

	CreatedAt time.Time

	UpdatedAt time.Time
}
