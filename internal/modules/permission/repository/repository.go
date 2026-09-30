package repository

import (
	"context"

	permissionmodel "github.com/dev-gopi/authhub/internal/modules/permission/model"

	"gorm.io/gorm"
)

type postgresRepository struct{}

func NewPostgresRepository() Interface {
	return &postgresRepository{}
}

func (r *postgresRepository) FindByAPILabels(
	ctx context.Context,
	tx *gorm.DB,
	labels []string,
) ([]permissionmodel.Permission, error) {
	var permissions []permissionmodel.Permission

	if len(labels) == 0 {
		return permissions, nil
	}

	err := tx.
		WithContext(ctx).
		Where(
			`api_label IN ?
			AND is_active = true
			AND is_deleted = false`,
			labels,
		).
		Find(&permissions).
		Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}
