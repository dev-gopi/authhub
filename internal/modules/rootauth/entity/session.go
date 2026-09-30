package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID uuid.UUID

	PlatformUserID uuid.UUID

	ExpiresAt time.Time

	IdleExpiresAt time.Time

	LastSeenAt time.Time
}
