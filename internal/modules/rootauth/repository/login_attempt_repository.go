package repository

import (
	"context"

	rootmodel "github.com/dev-gopi/authhub/internal/modules/rootauth/model"

	"gorm.io/gorm"
)

type loginAttemptRepository struct {
	db *gorm.DB
}

func NewLoginAttemptRepository(
	db *gorm.DB,
) LoginAttemptRepository {
	return &loginAttemptRepository{
		db: db,
	}
}

func (r *loginAttemptRepository) Create(
	ctx context.Context,
	attempt *rootmodel.PlatformLoginAttempt,
) error {
	return r.db.
		WithContext(ctx).
		Create(attempt).
		Error
}
