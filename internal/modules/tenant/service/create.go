package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	rootentity "github.com/dev-gopi/authhub/internal/modules/rootauth/entity"
	"github.com/dev-gopi/authhub/internal/modules/tenant/dto"
	tenantentity "github.com/dev-gopi/authhub/internal/modules/tenant/entity"
	tenantmodel "github.com/dev-gopi/authhub/internal/modules/tenant/model"
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

func (s *Service) Create(
	ctx context.Context,
	req dto.CreateTenantRequest,
	rootActor *rootentity.RootAuthContext,
) (*dto.CreateTenantResponse, error) {
	if rootActor == nil ||
		!rootActor.IsRootAdmin {
		return nil, errors.New(
			"root authorization required",
		)
	}

	existing, err :=
		s.tenants.FindByAPILabel(
			ctx,
			req.APILabel,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"check tenant api label: %w",
			err,
		)
	}

	if existing != nil {
		return nil, errors.New(
			"tenant api label already exists",
		)
	}

	now := time.Now().UTC()

	tenant := &tenantmodel.Tenant{
		BaseModel: sharedmodel.BaseModel{
			ID: uuid.New(),

			CreatedAt: now,
			UpdatedAt: now,

			IsActive: true,

			IsDeleted: false,
		},

		APILabel: req.APILabel,

		DisplayName: req.DisplayName,

		Status: tenantentity.
			StatusProvisioning.
			String(),

		DefaultLocale: "en",
	}

	// Important:
	//
	// This transaction creates ONLY the minimal tenant
	// record.
	//
	// We intentionally commit it before calling Vault.
	err = s.db.
		WithContext(ctx).
		Transaction(
			func(tx *gorm.DB) error {
				return s.tenants.Create(
					ctx,
					tx,
					tenant,
				)
			},
		)

	if err != nil {
		return nil, fmt.Errorf(
			"create provisioning tenant: %w",
			err,
		)
	}

	provisioned, err := s.provision(
		ctx,
		tenant,
		req,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"tenant provisioning failed: %w",
			err,
		)
	}
	_ = provisioned

	return &dto.CreateTenantResponse{
		ID: tenant.ID,

		APILabel: tenant.APILabel,

		DisplayName: tenant.DisplayName,

		Status: tenant.Status,
	}, nil
}
