package repository

import (
	"context"

	tenantmodel "github.com/dev-gopi/authhub/internal/modules/tenant/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Interface interface {
	Create(
		ctx context.Context,
		tx *gorm.DB,
		tenant *tenantmodel.Tenant,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*tenantmodel.Tenant, error)

	FindByAPILabel(
		ctx context.Context,
		apiLabel string,
	) (*tenantmodel.Tenant, error)

	UpdateStatus(
		ctx context.Context,
		tx *gorm.DB,
		tenantID uuid.UUID,
		status string,
	) error
}
