package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"github.com/google/uuid"
)

type TenantMemberRole struct {
	sharedmodel.BaseModel
	TenantMemberID uuid.UUID `gorm:"type:uuid;not null;index"`
	RoleID         uuid.UUID `gorm:"type:uuid;not null;index"`
}

func (TenantMemberRole) TableName() string { return "tenant_member_roles" }
