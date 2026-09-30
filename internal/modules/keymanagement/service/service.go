package service

import (
	"context"
	"fmt"

	keyconstants "github.com/dev-gopi/authhub/internal/modules/keymanagement/constants"
	keyentity "github.com/dev-gopi/authhub/internal/modules/keymanagement/entity"
	vaultinfra "github.com/dev-gopi/authhub/internal/platform/vault"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	vault *vaultinfra.Client

	transitMount string

	logger *zap.Logger
}

func New(
	vault *vaultinfra.Client,
	transitMount string,
	logger *zap.Logger,
) *Service {
	return &Service{
		vault: vault,

		transitMount: transitMount,

		logger: logger,
	}
}

func (s *Service) ProvisionTenantKeys(
	ctx context.Context,
	tenantID uuid.UUID,
) error {
	keySpecs := []vaultinfra.TransitKeySpec{
		{
			Name: keyentity.TenantKeyName(
				tenantID,
				keyconstants.KeyPurposeOIDCSigning,
			),

			Type: keyconstants.VaultKeyTypeEd25519,
		},
		{
			Name: keyentity.TenantKeyName(
				tenantID,
				keyconstants.KeyPurposeSAMLSigning,
			),

			Type: keyconstants.VaultKeyTypeRSA3072,
		},
		{
			Name: keyentity.TenantKeyName(
				tenantID,
				keyconstants.KeyPurposeWebhookSigning,
			),

			Type: keyconstants.VaultKeyTypeEd25519,
		},
		{
			Name: keyentity.TenantKeyName(
				tenantID,
				keyconstants.KeyPurposeRealtimeSigning,
			),

			Type: keyconstants.VaultKeyTypeEd25519,
		},
		{
			Name: keyentity.TenantKeyName(
				tenantID,
				keyconstants.KeyPurposeDataEncryption,
			),

			Type: keyconstants.VaultKeyTypeAES256GCM96,
		},
	}

	for _, spec := range keySpecs {
		if err := s.vault.EnsureTransitKey(
			ctx,
			s.transitMount,
			spec,
		); err != nil {
			return fmt.Errorf(
				"provision tenant vault key %q: %w",
				spec.Name,
				err,
			)
		}

		s.logger.Info(
			"tenant vault key ensured",
			zap.String(
				"tenant_id",
				tenantID.String(),
			),
			zap.String(
				"key_name",
				spec.Name,
			),
		)
	}

	return nil
}
