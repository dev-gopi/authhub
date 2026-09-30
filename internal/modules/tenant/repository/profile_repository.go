package repository

import (
	"context"

	tenantmodel "github.com/dev-gopi/authhub/internal/modules/tenant/model"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	Create(
		ctx context.Context,
		tx *gorm.DB,
		profile *tenantmodel.TenantProfile,
	) error
}

type profileRepository struct{}

func NewProfileRepository() ProfileRepository {
	return &profileRepository{}
}

func (r *profileRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	profile *tenantmodel.TenantProfile,
) error {
	return tx.
		WithContext(ctx).
		Create(profile).
		Error
}
