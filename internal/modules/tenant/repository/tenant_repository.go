package repository

import (
	"context"
	"errors"
	"time"

	tenantmodel "github.com/dev-gopi/authhub/internal/modules/tenant/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(
	db *gorm.DB,
) Interface {
	return &postgresRepository{
		db: db,
	}
}

func (r *postgresRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	tenant *tenantmodel.Tenant,
) error {
	return tx.
		WithContext(ctx).
		Create(tenant).
		Error
}

func (r *postgresRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*tenantmodel.Tenant, error) {
	var tenant tenantmodel.Tenant

	err := r.db.
		WithContext(ctx).
		Where(
			`id = ?
			AND is_deleted = false`,
			id,
		).
		First(&tenant).
		Error

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, nil
		}

		return nil, err
	}

	return &tenant, nil
}

func (r *postgresRepository) FindByAPILabel(
	ctx context.Context,
	apiLabel string,
) (*tenantmodel.Tenant, error) {
	var tenant tenantmodel.Tenant

	err := r.db.
		WithContext(ctx).
		Where(
			`api_label = ?
			AND is_deleted = false`,
			apiLabel,
		).
		First(&tenant).
		Error

	if err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			return nil, nil
		}

		return nil, err
	}

	return &tenant, nil
}

func (r *postgresRepository) UpdateStatus(
	ctx context.Context,
	tx *gorm.DB,
	tenantID uuid.UUID,
	status string,
) error {
	return tx.
		WithContext(ctx).
		Model(&tenantmodel.Tenant{}).
		Where(
			"id = ? AND is_deleted = false",
			tenantID,
		).
		Updates(map[string]any{
			"status":     status,
			"updated_at": time.Now().UTC(),
		}).
		Error
}
