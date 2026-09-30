package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"

	"github.com/google/uuid"
)

type AuthPolicy struct {
	sharedmodel.BaseModel

	UserPoolID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_pool_id"`

	UsernameLoginEnabled bool `gorm:"not null;default:true"`

	EmailLoginEnabled bool `gorm:"not null;default:false"`

	PasswordLoginEnabled bool `gorm:"not null;default:true"`

	PasskeyLoginEnabled bool `gorm:"not null;default:true"`

	RegistrationEnabled bool `gorm:"not null;default:true"`

	EmailVerificationRequired bool `gorm:"not null;default:true"`

	MFAPolicy string `gorm:"type:text;not null;default:optional"`

	TOTPEnabled bool `gorm:"not null;default:true"`

	SessionIdleSeconds int `gorm:"not null"`

	SessionAbsoluteSeconds int `gorm:"not null"`

	MaxActiveSessions *int

	RecentAuthSeconds int `gorm:"not null"`
}

func (AuthPolicy) TableName() string {
	return "user_pool_auth_policies"
}
