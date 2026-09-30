package model

import (
	"time"

	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"

	"github.com/google/uuid"
)

type PlatformLoginAttempt struct {
	sharedmodel.BaseModel

	PlatformUserID *uuid.UUID `gorm:"type:uuid;index"`

	IdentifierHash string `gorm:"type:char(64);not null;index"`

	IPAddress *string `gorm:"type:inet"`

	UserAgent *string `gorm:"type:text"`

	Success bool `gorm:"not null"`

	FailureReason *string `gorm:"type:text"`

	AttemptedAt time.Time `gorm:"type:timestamptz;not null;index"`
}

func (PlatformLoginAttempt) TableName() string {
	return "platform_login_attempts"
}
