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
)

func (s *Service) Create(ctx context.Context, req dto.CreateTenantRequest, rootActor *rootentity.RootAuthContext) (*dto.CreateTenantResponse, error) {
	if rootActor == nil || !rootActor.IsRootAdmin {
		return nil, tenantentity.ErrUnauthorized
	}

	tenant, err := s.tenants.FindByAPILabel(ctx, req.APILabel)
	if err != nil {
		return nil, fmt.Errorf("check tenant api label: %w", err)
	}
	if tenant != nil && tenant.Status != tenantentity.StatusProvisioning.String() {
		return nil, tenantentity.ErrTenantAlreadyExists
	}

	if tenant == nil {
		now := time.Now().UTC()
		tenant = &tenantmodel.Tenant{
			BaseModel: sharedmodel.NewBaseModelAt(now, &rootActor.PlatformUserID),
			APILabel:  req.APILabel, DisplayName: req.DisplayName,
			Status: tenantentity.StatusProvisioning.String(), DefaultLocale: "en",
		}
		if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return s.tenants.Create(ctx, tx, tenant) }); err != nil {
			return nil, fmt.Errorf("create provisioning tenant: %w", err)
		}
	}

	provisioned, err := s.provision(ctx, tenant, req, rootActor.PlatformUserID)
	if err != nil {
		s.recordProvisioningFailureAudit(ctx, tenant.ID, rootActor.PlatformUserID, err)
		if errors.Is(err, tenantentity.ErrPrimaryAdminConflict) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", tenantentity.ErrProvisioningFailed, err)
	}

	tenant.Status = tenantentity.StatusActive.String()
	return &dto.CreateTenantResponse{
		ID: tenant.ID, APILabel: tenant.APILabel, DisplayName: tenant.DisplayName, Status: tenant.Status,
		PrimaryAdmin:    dto.PrimaryAdminResponse{ID: provisioned.PrimaryAdminID, Username: req.PrimaryAdmin.Username, Email: req.PrimaryAdmin.Email},
		DefaultUserPool: dto.UserPoolResponse{ID: provisioned.UserPoolID, APILabel: "default"},
	}, nil
}
