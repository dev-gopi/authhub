package model

import (
	"time"

	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"

	"gorm.io/datatypes"
)

type Tenant struct {
	sharedmodel.BaseModel

	APILabel string `gorm:"type:text;not null;uniqueIndex" json:"api_label"`

	DisplayName string `gorm:"type:text;not null" json:"display_name"`

	Status string `gorm:"type:text;not null;index" json:"status"`

	DefaultLocale string `gorm:"type:text;not null;default:en" json:"default_locale"`

	Metadata datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"metadata"`

	SuspendedAt *time.Time `gorm:"type:timestamptz" json:"suspended_at,omitempty"`
}

func (Tenant) TableName() string {
	return "tenants"
}
