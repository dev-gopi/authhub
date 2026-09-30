package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
)

type Permission struct {
	sharedmodel.BaseModel

	APILabel string `gorm:"type:text;not null;uniqueIndex" json:"api_label"`

	DisplayName string `gorm:"type:text;not null" json:"display_name"`

	Description *string `gorm:"type:text" json:"description,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}
