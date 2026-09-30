package model

import (
	"time"

	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PlatformPassword struct {
	sharedmodel.BaseModel

	PlatformUserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	PasswordHash      string         `gorm:"not null"`
	PasswordAlgorithm string         `gorm:"not null"`
	PasswordParams    datatypes.JSON `gorm:"type:jsonb;not null"`
	PasswordVersion   int            `gorm:"not null"`

	ChangedAt time.Time `gorm:"not null"`
}

func (PlatformPassword) TableName() string {
	return "platform_user_passwords"
}
