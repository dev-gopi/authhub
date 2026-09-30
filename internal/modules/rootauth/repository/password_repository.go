package repository

import (
	"context"
	"errors"
	"time"

	rootmodel "github.com/dev-gopi/authhub/internal/modules/rootauth/model"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type passwordRepository struct {
	db *gorm.DB
}

func NewPasswordRepository(
	db *gorm.DB,
) PasswordRepository {
	return &passwordRepository{
		db: db,
	}
}

func (r *passwordRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	password *rootmodel.PlatformPassword,
) error {
	return tx.WithContext(ctx).Create(password).Error
}

func (r *passwordRepository) FindByPlatformUserID(
	ctx context.Context,
	platformUserID uuid.UUID,
) (*rootmodel.PlatformPassword, error) {
	var password rootmodel.PlatformPassword

	err := r.db.
		WithContext(ctx).
		Where(
			`platform_user_id = ?
			AND is_active = true
			AND is_deleted = false`,
			platformUserID,
		).
		First(&password).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &password, nil
}

func (r *passwordRepository) UpdateHash(
	ctx context.Context,
	tx *gorm.DB,
	platformUserID uuid.UUID,
	passwordHash string,
	passwordParams []byte,
) error {
	now := time.Now().UTC()

	return tx.
		WithContext(ctx).
		Model(&rootmodel.PlatformPassword{}).
		Where(
			`platform_user_id = ?
			AND is_active = true
			AND is_deleted = false`,
			platformUserID,
		).
		Updates(map[string]any{
			"password_hash":    passwordHash,
			"password_params":  datatypes.JSON(passwordParams),
			"password_version": gorm.Expr("password_version + 1"),
			"changed_at":       now,
			"updated_at":       now,
		}).
		Error
}
