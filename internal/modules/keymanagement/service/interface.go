package service

import (
	"context"

	"github.com/google/uuid"
)

type Interface interface {
	ProvisionTenantKeys(
		ctx context.Context,
		tenantID uuid.UUID,
	) error
}
