package repository

import (
	"context"
	auditmodel "github.com/dev-gopi/authhub/internal/modules/audit/model"
	"gorm.io/gorm"
)

type Interface interface {
	Create(context.Context, *gorm.DB, *auditmodel.Event) error
}

type postgresRepository struct{}

func NewPostgresRepository() Interface { return &postgresRepository{} }
func (r *postgresRepository) Create(ctx context.Context, tx *gorm.DB, event *auditmodel.Event) error {
	return tx.WithContext(ctx).Create(event).Error
}
