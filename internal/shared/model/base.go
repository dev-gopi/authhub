package model

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`

	UpdatedAt time.Time `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`

	DeletedAt *time.Time `gorm:"type:timestamptz;index" json:"deleted_at,omitempty"`

	IsActive bool `gorm:"not null;default:true" json:"is_active"`

	IsDeleted bool `gorm:"not null;default:false" json:"is_deleted"`
}
