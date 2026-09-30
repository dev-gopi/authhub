package repository

import (
	"context"
	"time"

	platformmodel "github.com/dev-gopi/authhub/internal/modules/platformuser/model"
	rootmodel "github.com/dev-gopi/authhub/internal/modules/rootauth/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlatformUserRepository interface {
	FindByIdentifier(
		ctx context.Context,
		identifier string,
	) (*platformmodel.PlatformUser, error)

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*platformmodel.PlatformUser, error)
}

type PasswordRepository interface {
	FindByPlatformUserID(
		ctx context.Context,
		platformUserID uuid.UUID,
	) (*rootmodel.PlatformPassword, error)

	UpdateHash(
		ctx context.Context,
		tx *gorm.DB,
		platformUserID uuid.UUID,
		passwordHash string,
		passwordParams []byte,
	) error
}

type SessionRepository interface {
	Create(
		ctx context.Context,
		session *rootmodel.PlatformSession,
	) error

	FindActiveByTokenHash(
		ctx context.Context,
		tokenHash string,
	) (*rootmodel.PlatformSession, error)

	Touch(
		ctx context.Context,
		sessionID uuid.UUID,
		lastSeenAt time.Time,
		idleExpiresAt time.Time,
	) error

	Revoke(
		ctx context.Context,
		sessionID uuid.UUID,
		reason string,
	) error

	RevokeAllByPlatformUser(
		ctx context.Context,
		platformUserID uuid.UUID,
		reason string,
	) error
}
