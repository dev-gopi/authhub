package model

import (
	"time"

	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
)

// PlatformUser represents a control-plane identity.
//
// Root Admin and tenant administrators are both platform users,
// but root authorization is controlled exclusively through
// IsRootAdmin.
//
// Tenant membership and tenant roles are stored separately.
type PlatformUser struct {
	sharedmodel.BaseModel

	Username string `gorm:"type:text;not null;uniqueIndex:idx_platform_users_username" json:"username"`

	Email string `gorm:"type:text;not null;uniqueIndex:idx_platform_users_email" json:"email"`

	EmailVerifiedAt *time.Time `gorm:"type:timestamptz" json:"email_verified_at,omitempty"`

	DisplayName string `gorm:"type:text" json:"display_name"`

	Status string `gorm:"type:text;not null;index:idx_platform_users_status" json:"status"`

	IsRootAdmin bool `gorm:"not null;default:false;index:idx_platform_users_root_admin" json:"is_root_admin"`

	CredentialVersion int64 `gorm:"not null;default:1" json:"credential_version"`

	MustChangePassword bool `gorm:"not null;default:false" json:"must_change_password"`

	LastLoginAt *time.Time `gorm:"type:timestamptz" json:"last_login_at,omitempty"`
}

func (PlatformUser) TableName() string {
	return "platform_users"
}
