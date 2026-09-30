package model

import (
	"time"

	"github.com/google/uuid"

	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
)

type PlatformSession struct {
	sharedmodel.BaseModel

	PlatformUserID uuid.UUID `gorm:"type:uuid;not null;index"`

	SessionTokenHash string `gorm:"type:char(64);not null;uniqueIndex"`

	CredentialVersion int64 `gorm:"not null"`

	IPAddress *string

	UserAgent *string

	ExpiresAt time.Time `gorm:"not null"`

	IdleExpiresAt time.Time `gorm:"not null"`

	LastSeenAt time.Time `gorm:"not null"`

	RevokedAt *time.Time

	RevokedReason *string
}

func (PlatformSession) TableName() string {
	return "platform_user_sessions"
}
