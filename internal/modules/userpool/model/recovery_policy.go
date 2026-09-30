package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"

	"github.com/google/uuid"
)

type RecoveryPolicy struct {
	sharedmodel.BaseModel

	UserPoolID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`

	EmailLinkEnabled bool `gorm:"not null;default:true"`

	EmailOTPEnabled bool `gorm:"not null;default:true"`

	SMSOTPEnabled bool `gorm:"not null;default:false"`

	SavedRecoveryCodeEnabled bool `gorm:"not null;default:true"`

	ExistingTOTPEnabled bool `gorm:"not null;default:true"`

	ExistingPasskeyEnabled bool `gorm:"not null;default:true"`

	AdminAssistedEnabled bool `gorm:"not null;default:true"`

	SecurityQuestionEnabled bool `gorm:"not null;default:false"`

	MinimumProofs int `gorm:"not null;default:1"`

	RequireIndependentChannels bool `gorm:"not null;default:false"`

	PostRecoverySessionAction string `gorm:"type:text;not null;default:revoke_all"`
}

func (RecoveryPolicy) TableName() string {
	return "recovery_policies"
}
