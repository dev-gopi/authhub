package repository

import (
	"context"
	"errors"
	"time"

	rootmodel "github.com/dev-gopi/authhub/internal/modules/rootauth/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(
	db *gorm.DB,
) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) Create(
	ctx context.Context,
	session *rootmodel.PlatformSession,
) error {
	return r.db.
		WithContext(ctx).
		Create(session).
		Error
}

func (r *sessionRepository) FindActiveByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*rootmodel.PlatformSession, error) {
	var session rootmodel.PlatformSession

	err := r.db.
		WithContext(ctx).
		Where(
			`session_token_hash = ?
			AND revoked_at IS NULL
			AND is_active = true
			AND is_deleted = false`,
			tokenHash,
		).
		First(&session).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) Touch(
	ctx context.Context,
	sessionID uuid.UUID,
	lastSeenAt time.Time,
	idleExpiresAt time.Time,
) error {
	return r.db.
		WithContext(ctx).
		Model(&rootmodel.PlatformSession{}).
		Where(
			"id = ? AND revoked_at IS NULL",
			sessionID,
		).
		Updates(map[string]any{
			"last_seen_at":    lastSeenAt,
			"idle_expires_at": idleExpiresAt,
			"updated_at":      time.Now().UTC(),
		}).
		Error
}

func (r *sessionRepository) Revoke(
	ctx context.Context,
	sessionID uuid.UUID,
	reason string,
) error {
	now := time.Now().UTC()

	return r.db.
		WithContext(ctx).
		Model(&rootmodel.PlatformSession{}).
		Where(
			"id = ? AND revoked_at IS NULL",
			sessionID,
		).
		Updates(map[string]any{
			"revoked_at":     now,
			"revoked_reason": reason,
			"is_active":      false,
			"updated_at":     now,
		}).
		Error
}

func (r *sessionRepository) RevokeAllByPlatformUser(
	ctx context.Context,
	platformUserID uuid.UUID,
	reason string,
) error {
	now := time.Now().UTC()

	return r.db.
		WithContext(ctx).
		Model(&rootmodel.PlatformSession{}).
		Where(
			`platform_user_id = ?
			AND revoked_at IS NULL`,
			platformUserID,
		).
		Updates(map[string]any{
			"revoked_at":     now,
			"revoked_reason": reason,
			"is_active":      false,
			"updated_at":     now,
		}).
		Error
}
