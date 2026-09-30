package repository

import (
	"context"
	"errors"
	"time"

	platformmodel "github.com/dev-gopi/authhub/internal/modules/platformuser/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type platformUserRepository struct {
	db *gorm.DB
}

func NewPlatformUserRepository(
	db *gorm.DB,
) PlatformUserRepository {
	return &platformUserRepository{
		db: db,
	}
}

func (r *platformUserRepository) FindByIdentifier(
	ctx context.Context,
	identifier string,
) (*platformmodel.PlatformUser, error) {
	var user platformmodel.PlatformUser

	err := r.db.
		WithContext(ctx).
		Where(
			`(
				LOWER(username) = ?
				OR LOWER(email) = ?
			)
			AND is_deleted = false`,
			identifier,
			identifier,
		).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *platformUserRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*platformmodel.PlatformUser, error) {
	var user platformmodel.PlatformUser

	err := r.db.
		WithContext(ctx).
		Where(
			"id = ? AND is_deleted = false",
			id,
		).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *platformUserRepository) UpdateLastLogin(
	ctx context.Context,
	id uuid.UUID,
	at time.Time,
) error {
	return r.db.
		WithContext(ctx).
		Model(&platformmodel.PlatformUser{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"last_login_at": at,
			"updated_at":    at,
		}).
		Error
}
