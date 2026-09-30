package service

import (
	"context"
	"github.com/google/uuid"
)

type Interface interface {
	ProvisionTenantKeys(ctx context.Context, tenantID uuid.UUID) error
	EncryptTenantData(ctx context.Context, tenantID uuid.UUID, plaintext []byte) (string, error)
}
