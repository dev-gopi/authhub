package repository

import (
	"context"
	"errors"

	membermodel "github.com/dev-gopi/authhub/internal/modules/tenantmember/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Interface interface {
	Create(context.Context, *gorm.DB, *membermodel.TenantMember) error
	FindPrimaryAdmin(context.Context, *gorm.DB, uuid.UUID) (*membermodel.TenantMember, error)
}

type postgresRepository struct{}

func NewPostgresRepository() Interface { return &postgresRepository{} }
func (r *postgresRepository) Create(ctx context.Context, tx *gorm.DB, m *membermodel.TenantMember) error {
	return tx.WithContext(ctx).Create(m).Error
}
func (r *postgresRepository) FindPrimaryAdmin(ctx context.Context, tx *gorm.DB, tenantID uuid.UUID) (*membermodel.TenantMember, error) {
	var m membermodel.TenantMember
	err := tx.WithContext(ctx).Where("tenant_id = ? AND member_class = 'primary_admin' AND status = 'active' AND is_active = true AND is_deleted = false", tenantID).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}
