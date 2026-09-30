package model

import (
	"time"

	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Event struct {
	sharedmodel.BaseModel
	TenantID      *uuid.UUID     `gorm:"type:uuid;index"`
	EventType     string         `gorm:"type:text;not null;index"`
	AggregateType string         `gorm:"type:text;not null"`
	AggregateID   *uuid.UUID     `gorm:"type:uuid"`
	Payload       datatypes.JSON `gorm:"type:jsonb;not null"`
	Status        string         `gorm:"type:text;not null;index"`
	Attempts      int            `gorm:"not null;default:0"`
	AvailableAt   time.Time      `gorm:"type:timestamptz;not null"`
	PublishedAt   *time.Time     `gorm:"type:timestamptz"`
	LastError     *string        `gorm:"type:text"`
}

func (Event) TableName() string { return "outbox_events" }
