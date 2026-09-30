package repository

import (
	"context"
	"errors"

	rolemodel "github.com/dev-gopi/authhub/internal/modules/role/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Interface interface {
	EnsureRole(context.Context, *gorm.DB, *rolemodel.Role) (*rolemodel.Role, error)
	EnsureRolePermission(context.Context, *gorm.DB, *rolemodel.RolePermission) error
	EnsureTenantMemberRole(context.Context, *gorm.DB, *rolemodel.TenantMemberRole) error
	FindByAPILabel(context.Context, *gorm.DB, uuid.UUID, string) (*rolemodel.Role, error)
	ListDefaultRoleSnapshots(context.Context, *gorm.DB, string) ([]rolemodel.DefaultRoleSnapshot, error)
}

type postgresRepository struct{}

func NewPostgresRepository() Interface { return &postgresRepository{} }

func (r *postgresRepository) EnsureRole(ctx context.Context, tx *gorm.DB, role *rolemodel.Role) (*rolemodel.Role, error) {
	if err := tx.WithContext(ctx).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "tenant_id"}, {Name: "api_label"}}, DoNothing: true}).Create(role).Error; err != nil {
		return nil, err
	}
	var stored rolemodel.Role
	if err := tx.WithContext(ctx).Where("tenant_id = ? AND api_label = ? AND is_deleted = false", role.TenantID, role.APILabel).First(&stored).Error; err != nil {
		return nil, err
	}
	return &stored, nil
}

func (r *postgresRepository) EnsureRolePermission(ctx context.Context, tx *gorm.DB, rp *rolemodel.RolePermission) error {
	return tx.WithContext(ctx).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "role_id"}, {Name: "permission_id"}}, DoNothing: true}).Create(rp).Error
}

func (r *postgresRepository) EnsureTenantMemberRole(ctx context.Context, tx *gorm.DB, mr *rolemodel.TenantMemberRole) error {
	return tx.WithContext(ctx).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "tenant_member_id"}, {Name: "role_id"}}, DoNothing: true}).Create(mr).Error
}

func (r *postgresRepository) FindByAPILabel(ctx context.Context, tx *gorm.DB, tenantID uuid.UUID, apiLabel string) (*rolemodel.Role, error) {
	var role rolemodel.Role
	err := tx.WithContext(ctx).Where("tenant_id = ? AND api_label = ? AND is_deleted = false", tenantID, apiLabel).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *postgresRepository) ListDefaultRoleSnapshots(ctx context.Context, tx *gorm.DB, platformVersion string) ([]rolemodel.DefaultRoleSnapshot, error) {
	var snapshots []rolemodel.DefaultRoleSnapshot
	err := tx.WithContext(ctx).
		Where("platform_version = ? AND is_active = true AND is_deleted = false", platformVersion).
		Order("role_api_label ASC").
		Find(&snapshots).Error
	if err != nil {
		return nil, err
	}
	return snapshots, nil
}
