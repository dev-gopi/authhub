package repository

import (
	"context"
	outboxmodel "github.com/dev-gopi/authhub/internal/modules/outbox/model"
	"gorm.io/gorm"
)

type Interface interface {
	Create(context.Context, *gorm.DB, *outboxmodel.Event) error
}

type postgresRepository struct{}

func NewPostgresRepository() Interface { return &postgresRepository{} }
func (r *postgresRepository) Create(ctx context.Context, tx *gorm.DB, event *outboxmodel.Event) error {
	return tx.WithContext(ctx).Create(event).Error
}
