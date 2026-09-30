package repository

import (
	"context"

	permissionmodel "github.com/dev-gopi/authhub/internal/modules/permission/model"

	"gorm.io/gorm"
)

type Interface interface {
	FindByAPILabels(
		ctx context.Context,
		tx *gorm.DB,
		labels []string,
	) ([]permissionmodel.Permission, error)
}
