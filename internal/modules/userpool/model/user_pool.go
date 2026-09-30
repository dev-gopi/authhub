package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"

	"github.com/google/uuid"
)

type UserPool struct {
	sharedmodel.BaseModel

	TenantID uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`

	APILabel string `gorm:"type:text;not null" json:"api_label"`

	DisplayName string `gorm:"type:text;not null" json:"display_name"`

	Status string `gorm:"type:text;not null" json:"status"`

	IssuerURI string `gorm:"type:text;not null" json:"issuer_uri"`

	DefaultLocale string `gorm:"type:text;not null;default:en" json:"default_locale"`
}

func (UserPool) TableName() string {
	return "user_pools"
}
