package model

import (
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Event struct {
	sharedmodel.BaseModel
	TenantID            *uuid.UUID     `gorm:"type:uuid;index"`
	ActorPlatformUserID *uuid.UUID     `gorm:"type:uuid;index"`
	Action              string         `gorm:"type:text;not null;index"`
	ResourceType        string         `gorm:"type:text;not null"`
	ResourceID          *uuid.UUID     `gorm:"type:uuid"`
	Success             bool           `gorm:"not null"`
	Metadata            datatypes.JSON `gorm:"type:jsonb;not null"`
}

func (Event) TableName() string { return "audit_events" }
