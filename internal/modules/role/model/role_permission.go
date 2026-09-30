package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"github.com/google/uuid"
)

type RolePermission struct {
	sharedmodel.BaseModel
	RoleID       uuid.UUID `gorm:"type:uuid;not null;index"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null;index"`
}

func (RolePermission) TableName() string { return "role_permissions" }
