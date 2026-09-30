package model

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`

	UpdatedAt time.Time `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`

	CreatedBy *uuid.UUID `gorm:"type:uuid;index" json:"created_by,omitempty"`

	UpdatedBy *uuid.UUID `gorm:"type:uuid;index" json:"updated_by,omitempty"`

	DeletedAt *time.Time `gorm:"type:timestamptz;index" json:"deleted_at,omitempty"`

	IsActive bool `gorm:"not null;default:true" json:"is_active"`

	IsDeleted bool `gorm:"not null;default:false" json:"is_deleted"`
}

func NewBaseModel(actorID *uuid.UUID) BaseModel {
	return NewBaseModelAt(time.Now().UTC(), actorID)
}

func NewBaseModelAt(now time.Time, actorID *uuid.UUID) BaseModel {
	return BaseModel{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: actorID,
		UpdatedBy: actorID,
		IsActive:  true,
		IsDeleted: false,
	}
}
