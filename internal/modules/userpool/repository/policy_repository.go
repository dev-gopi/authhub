package repository

import (
	"context"

	userpoolmodel "github.com/dev-gopi/authhub/internal/modules/userpool/model"

	"gorm.io/gorm"
)

type postgresRepository struct{}

func NewPostgresRepository() Interface {
	return &postgresRepository{}
}

func (r *postgresRepository) CreateUserPool(
	ctx context.Context,
	tx *gorm.DB,
	pool *userpoolmodel.UserPool,
) error {
	return tx.
		WithContext(ctx).
		Create(pool).
		Error
}

func (r *postgresRepository) CreateAuthPolicy(
	ctx context.Context,
	tx *gorm.DB,
	policy *userpoolmodel.AuthPolicy,
) error {
	return tx.
		WithContext(ctx).
		Create(policy).
		Error
}

func (r *postgresRepository) CreatePasswordPolicy(
	ctx context.Context,
	tx *gorm.DB,
	policy *userpoolmodel.PasswordPolicy,
) error {
	return tx.
		WithContext(ctx).
		Create(policy).
		Error
}

func (r *postgresRepository) CreateRecoveryPolicy(
	ctx context.Context,
	tx *gorm.DB,
	policy *userpoolmodel.RecoveryPolicy,
) error {
	return tx.
		WithContext(ctx).
		Create(policy).
		Error
}
