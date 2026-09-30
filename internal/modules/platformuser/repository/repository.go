package repository

import (
	"context"
	"errors"
	"strings"

	platformmodel "github.com/dev-gopi/authhub/internal/modules/platformuser/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Interface interface {
	Create(context.Context, *gorm.DB, *platformmodel.PlatformUser) error
	FindByUsernameOrEmail(context.Context, *gorm.DB, string, string) (*platformmodel.PlatformUser, error)
	FindByID(context.Context, *gorm.DB, uuid.UUID) (*platformmodel.PlatformUser, error)
}

type postgresRepository struct{}

func NewPostgresRepository() Interface { return &postgresRepository{} }

func (r *postgresRepository) Create(ctx context.Context, tx *gorm.DB, user *platformmodel.PlatformUser) error {
	return tx.WithContext(ctx).Create(user).Error
}

func (r *postgresRepository) FindByUsernameOrEmail(ctx context.Context, tx *gorm.DB, username, email string) (*platformmodel.PlatformUser, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	email = strings.ToLower(strings.TrimSpace(email))
	var user platformmodel.PlatformUser
	err := tx.WithContext(ctx).Where("(LOWER(username) = ? OR LOWER(email) = ?) AND is_deleted = false", username, email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepository) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*platformmodel.PlatformUser, error) {
	var user platformmodel.PlatformUser
	err := tx.WithContext(ctx).Where("id = ? AND is_deleted = false", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
