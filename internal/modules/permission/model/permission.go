package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"

	"gorm.io/datatypes"
)

type Permission struct {
	sharedmodel.BaseModel

	APILabel string `gorm:"type:text;not null;uniqueIndex" json:"api_label"`

	DisplayName string `gorm:"type:text;not null" json:"display_name"`

	Description *string `gorm:"type:text" json:"description,omitempty"`

	Metadata datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"metadata"`
}

func (Permission) TableName() string {
	return "permissions"
}
