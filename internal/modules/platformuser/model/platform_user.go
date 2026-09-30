package model

import (
	"time"

	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
)

// PlatformUser represents an identity that exists in the AuthHub
// control plane.
//
// A PlatformUser may be:
//   - Root Admin
//   - Primary Tenant Admin
//   - Delegated Tenant Admin
//
// Tenant membership is stored separately in tenant_members.
//
// Root privileges must never be inferred from tenant membership.
// They are explicitly controlled by IsRootAdmin.
type PlatformUser struct {
	sharedmodel.BaseModel

	Username string `gorm:"type:text;not null;uniqueIndex:idx_platform_users_username" json:"username"`

	Email string `gorm:"type:text;not null;uniqueIndex:idx_platform_users_email" json:"email"`

	EmailVerifiedAt *time.Time `gorm:"type:timestamptz" json:"email_verified_at,omitempty"`

	DisplayName string `gorm:"type:text" json:"display_name"`

	Status string `gorm:"type:text;not null;index:idx_platform_users_status" json:"status"`

	// IsRootAdmin explicitly identifies a platform user that is
	// permitted to operate the AuthHub root control plane.
	//
	// Authentication alone is not enough to access root resources.
	// Root authorization middleware must additionally require this
	// value to be true.
	IsRootAdmin bool `gorm:"not null;default:false;index:idx_platform_users_root_admin" json:"is_root_admin"`

	// CredentialVersion invalidates existing sessions whenever
	// security-sensitive credentials are changed.
	//
	// Example:
	//
	// user credential_version = 5
	// session credential_version = 4
	//
	// => session must be rejected.
	CredentialVersion int64 `gorm:"not null;default:1" json:"credential_version"`

	LastLoginAt *time.Time `gorm:"type:timestamptz" json:"last_login_at,omitempty"`
}

// TableName explicitly defines the PostgreSQL table used by GORM.
func (PlatformUser) TableName() string {
	return "platform_users"
}
