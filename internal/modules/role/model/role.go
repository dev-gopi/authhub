package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Role struct {
	sharedmodel.BaseModel

	TenantID *uuid.UUID `gorm:"type:uuid;index" json:"tenant_id,omitempty"`

	APILabel string `gorm:"type:text;not null" json:"api_label"`

	DisplayName string `gorm:"type:text;not null" json:"display_name"`

	Description *string `gorm:"type:text" json:"description,omitempty"`

	Metadata datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"metadata"`

	IsDefault bool `gorm:"not null;default:false" json:"is_default"`

	ProtectedFromDelete bool `gorm:"not null;default:false" json:"protected_from_delete"`
}

func (Role) TableName() string {
	return "roles"
}
