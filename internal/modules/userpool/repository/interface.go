package repository

import (
	"context"

	userpoolmodel "github.com/dev-gopi/authhub/internal/modules/userpool/model"

	"gorm.io/gorm"
)

type Interface interface {
	CreateUserPool(
		ctx context.Context,
		tx *gorm.DB,
		pool *userpoolmodel.UserPool,
	) error

	CreateAuthPolicy(
		ctx context.Context,
		tx *gorm.DB,
		policy *userpoolmodel.AuthPolicy,
	) error

	CreatePasswordPolicy(
		ctx context.Context,
		tx *gorm.DB,
		policy *userpoolmodel.PasswordPolicy,
	) error

	CreateRecoveryPolicy(
		ctx context.Context,
		tx *gorm.DB,
		policy *userpoolmodel.RecoveryPolicy,
	) error
}
