package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"gorm.io/datatypes"
)

type DefaultRoleSnapshot struct {
	sharedmodel.BaseModel

	RoleAPILabel string `gorm:"type:text;not null" json:"role_api_label"`

	PlatformVersion string `gorm:"type:text;not null" json:"platform_version"`

	Definition datatypes.JSON `gorm:"type:jsonb;not null" json:"definition"`

	DefinitionSHA256 string `gorm:"type:text;not null" json:"definition_sha256"`
}

func (DefaultRoleSnapshot) TableName() string {
	return "default_role_snapshots"
}
