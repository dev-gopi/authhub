package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"github.com/google/uuid"
)

type TenantMember struct {
	sharedmodel.BaseModel
	TenantID       uuid.UUID `gorm:"type:uuid;not null;index"`
	PlatformUserID uuid.UUID `gorm:"type:uuid;not null;index"`
	MemberClass    string    `gorm:"type:text;not null"`
	Status         string    `gorm:"type:text;not null"`
	ProtectedAdmin bool      `gorm:"not null;default:false"`
}

func (TenantMember) TableName() string { return "tenant_members" }
