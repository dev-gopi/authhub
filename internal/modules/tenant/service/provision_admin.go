package service

import (
	"context"
	"fmt"
	"time"

	platformmodel "github.com/dev-gopi/authhub/internal/modules/platformuser/model"
	rolemodel "github.com/dev-gopi/authhub/internal/modules/role/model"
	rootmodel "github.com/dev-gopi/authhub/internal/modules/rootauth/model"
	"github.com/dev-gopi/authhub/internal/modules/tenant/dto"
	tenantentity "github.com/dev-gopi/authhub/internal/modules/tenant/entity"
	membermodel "github.com/dev-gopi/authhub/internal/modules/tenantmember/model"
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type primaryAdminProvisionResult struct {
	PlatformUserID uuid.UUID
	TenantMemberID uuid.UUID
}

func (s *Service) provisionPrimaryAdmin(
	ctx context.Context,
	tx *gorm.DB,
	tenantID uuid.UUID,
	primaryAdminRoleID uuid.UUID,
	req dto.PrimaryAdminRequest,
	credential *preparedTemporaryCredential,
	actorID uuid.UUID,
	now time.Time,
) (*primaryAdminProvisionResult, error) {
	existingUser, err := s.platformUsers.FindByUsernameOrEmail(ctx, tx, req.Username, req.Email)
	if err != nil {
		return nil, fmt.Errorf("check primary admin identity: %w", err)
	}
	if existingUser != nil {
		return nil, tenantentity.ErrPrimaryAdminConflict
	}

	user := &platformmodel.PlatformUser{
		BaseModel: sharedmodel.NewBaseModelAt(now, &actorID),
		Username:  req.Username, Email: req.Email, EmailVerifiedAt: nil, DisplayName: req.DisplayName,
		Status: "active", IsRootAdmin: false, CredentialVersion: 1, MustChangePassword: true,
	}
	if err := s.platformUsers.Create(ctx, tx, user); err != nil {
		return nil, fmt.Errorf("create primary admin platform user: %w", err)
	}

	member := &membermodel.TenantMember{
		BaseModel: sharedmodel.NewBaseModelAt(now, &actorID),
		TenantID:  tenantID, PlatformUserID: user.ID, MemberClass: "primary_admin", Status: "active", ProtectedAdmin: true,
	}
	if err := s.tenantMembers.Create(ctx, tx, member); err != nil {
		return nil, fmt.Errorf("create primary admin tenant membership: %w", err)
	}

	memberRole := &rolemodel.TenantMemberRole{
		BaseModel:      sharedmodel.NewBaseModelAt(now, &actorID),
		TenantMemberID: member.ID, RoleID: primaryAdminRoleID,
	}
	if err := s.roles.EnsureTenantMemberRole(ctx, tx, memberRole); err != nil {
		return nil, fmt.Errorf("assign primary_admin role: %w", err)
	}

	password := &rootmodel.PlatformPassword{
		BaseModel:      sharedmodel.NewBaseModelAt(now, &actorID),
		PlatformUserID: user.ID, PasswordHash: credential.PasswordHash, PasswordAlgorithm: "argon2id",
		PasswordParams: datatypes.JSON(credential.PasswordParams), PasswordVersion: 1, ChangedAt: now,
	}
	if err := s.passwords.Create(ctx, tx, password); err != nil {
		return nil, fmt.Errorf("persist primary admin temporary password hash: %w", err)
	}

	return &primaryAdminProvisionResult{PlatformUserID: user.ID, TenantMemberID: member.ID}, nil
}
