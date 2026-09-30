package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PasswordPolicy struct {
	sharedmodel.BaseModel

	UserPoolID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	MinLength int `gorm:"not null"`

	MaxLength int `gorm:"not null"`

	PreventCommonPasswords bool `gorm:"not null;default:true"`

	BreachedPasswordCheckEnabled bool `gorm:"not null;default:false"`

	PasswordHistoryCount int `gorm:"not null;default:0"`

	ResetTokenTTLSeconds int `gorm:"not null"`

	LockoutPolicy datatypes.JSON `gorm:"type:jsonb;not null"`
}

func (PasswordPolicy) TableName() string {
	return "password_policies"
}
