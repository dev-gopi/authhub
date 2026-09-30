package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dev-gopi/authhub/internal/modules/tenant/dto"
	tenantmodel "github.com/dev-gopi/authhub/internal/modules/tenant/model"

	userpoolconstants "github.com/dev-gopi/authhub/internal/modules/userpool/constants"
	userpoolmodel "github.com/dev-gopi/authhub/internal/modules/userpool/model"

	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (s *Service) provision(
	ctx context.Context,
	tenant *tenantmodel.Tenant,
	req dto.CreateTenantRequest,
	actorID uuid.UUID,
) (*provisionResult, error) {
	/*
		STEP 1:
		Vault runs before the second PostgreSQL transaction.

		If this fails, tenant remains PROVISIONING.
	*/
	if err := s.keyManagement.ProvisionTenantKeys(
		ctx,
		tenant.ID,
	); err != nil {
		return nil, fmt.Errorf(
			"provision tenant keys: %w",
			err,
		)
	}

	credential, err := s.prepareTemporaryCredential(ctx, tenant.ID)
	if err != nil {
		return nil, fmt.Errorf("prepare primary admin temporary credential: %w", err)
	}

	result := &provisionResult{}

	/*
		STEP 2:
		Create relational provisioning state atomically.
	*/
	err = s.db.
		WithContext(ctx).
		Transaction(
			func(tx *gorm.DB) error {
				now := time.Now().UTC()

				profile := buildTenantProfile(
					tenant.ID,
					req.Profile,
					now,
					&actorID,
				)

				if err := s.profiles.Create(
					ctx,
					tx,
					profile,
				); err != nil {
					return fmt.Errorf(
						"create tenant profile: %w",
						err,
					)
				}

				pool := buildDefaultUserPool(
					tenant,
					s.issuerBaseURL,
					now,
					&actorID,
				)

				if err := s.userPools.CreateUserPool(
					ctx,
					tx,
					pool,
				); err != nil {
					return fmt.Errorf(
						"create default user pool: %w",
						err,
					)
				}

				result.UserPoolID = pool.ID

				authPolicy :=
					buildDefaultAuthPolicy(
						pool.ID,
						now,
						&actorID,
					)

				if err := s.userPools.CreateAuthPolicy(
					ctx,
					tx,
					authPolicy,
				); err != nil {
					return fmt.Errorf(
						"create default auth policy: %w",
						err,
					)
				}

				passwordPolicy :=
					buildDefaultPasswordPolicy(
						pool.ID,
						now,
						&actorID,
					)

				if err := s.userPools.CreatePasswordPolicy(
					ctx,
					tx,
					passwordPolicy,
				); err != nil {
					return fmt.Errorf(
						"create default password policy: %w",
						err,
					)
				}

				recoveryPolicy :=
					buildDefaultRecoveryPolicy(
						pool.ID,
						now,
						&actorID,
					)

				if err := s.userPools.CreateRecoveryPolicy(
					ctx,
					tx,
					recoveryPolicy,
				); err != nil {
					return fmt.Errorf(
						"create default recovery policy: %w",
						err,
					)
				}

				roleResult, err := s.provisionDefaultRoles(ctx, tx, tenant.ID, actorID, now)
				if err != nil {
					return fmt.Errorf("provision default roles: %w", err)
				}
				result.PrimaryAdminRoleID = roleResult.PrimaryAdminRoleID

				adminResult, err := s.provisionPrimaryAdmin(
					ctx, tx, tenant.ID, roleResult.PrimaryAdminRoleID, req.PrimaryAdmin, credential, actorID, now,
				)
				if err != nil {
					return fmt.Errorf("provision primary tenant admin: %w", err)
				}
				result.PrimaryAdminID = adminResult.PlatformUserID

				if err := s.createProvisioningEvents(
					ctx, tx, tenant.ID, actorID, adminResult.PlatformUserID, req.PrimaryAdmin.Email, credential.EncryptedPlaintext, now,
				); err != nil {
					return err
				}

				if err := s.tenants.UpdateStatus(ctx, tx, tenant.ID, "active"); err != nil {
					return fmt.Errorf("activate tenant: %w", err)
				}

				return nil
			},
		)

	if err != nil {
		return nil, fmt.Errorf(
			"provision tenant relational state: %w",
			err,
		)
	}

	return result, nil
}

func buildTenantProfile(
	tenantID uuid.UUID,
	req dto.TenantProfileRequest,
	now time.Time,
	actorID *uuid.UUID,
) *tenantmodel.TenantProfile {
	return &tenantmodel.TenantProfile{
		BaseModel: sharedmodel.NewBaseModelAt(now, actorID),

		TenantID: tenantID,

		LegalName: optionalString(
			req.LegalName,
		),

		WebsiteURL: optionalString(
			req.WebsiteURL,
		),

		SupportEmail: optionalString(
			req.SupportEmail,
		),

		SupportPhone: optionalString(
			req.SupportPhone,
		),

		SecurityContactEmail: optionalString(
			req.SecurityContactEmail,
		),

		Locale: optionalString(
			req.Locale,
		),

		Timezone: optionalString(
			req.Timezone,
		),

		CountryCode: optionalString(
			req.CountryCode,
		),

		Region: optionalString(
			req.Region,
		),

		PrivacyPolicyURL: optionalString(
			req.PrivacyPolicyURL,
		),

		TermsURL: optionalString(
			req.TermsURL,
		),

		Metadata: datatypes.JSON(
			[]byte(`{}`),
		),
	}
}

func buildDefaultUserPool(
	tenant *tenantmodel.Tenant,
	issuerBaseURL string,
	now time.Time,
	actorID *uuid.UUID,
) *userpoolmodel.UserPool {
	return &userpoolmodel.UserPool{
		BaseModel: sharedmodel.NewBaseModelAt(now, actorID),

		TenantID: tenant.ID,

		APILabel: userpoolconstants.
			DefaultPoolAPILabel,

		DisplayName: userpoolconstants.
			DefaultPoolDisplayName,

		Status: userpoolconstants.
			DefaultPoolStatus,

		IssuerURI: buildIssuerURI(
			issuerBaseURL,
			tenant.APILabel,
			userpoolconstants.DefaultPoolAPILabel,
		),

		DefaultLocale: userpoolconstants.
			DefaultLocale,
	}
}

func buildDefaultAuthPolicy(
	userPoolID uuid.UUID,
	now time.Time,
	actorID *uuid.UUID,
) *userpoolmodel.AuthPolicy {
	return &userpoolmodel.AuthPolicy{
		BaseModel: sharedmodel.NewBaseModelAt(now, actorID),

		UserPoolID: userPoolID,

		UsernameLoginEnabled: true,

		EmailLoginEnabled: false,

		PasswordLoginEnabled: true,

		PasskeyLoginEnabled: true,

		RegistrationEnabled: true,

		EmailVerificationRequired: true,

		MFAPolicy: userpoolconstants.
			DefaultMFAPolicy,

		TOTPEnabled: true,

		SessionIdleSeconds: userpoolconstants.
			DefaultSessionIdleSeconds,

		SessionAbsoluteSeconds: userpoolconstants.
			DefaultSessionAbsoluteSeconds,

		MaxActiveSessions: nil,

		RecentAuthSeconds: userpoolconstants.
			DefaultRecentAuthSeconds,
	}
}

func buildDefaultPasswordPolicy(
	userPoolID uuid.UUID,
	now time.Time,
	actorID *uuid.UUID,
) *userpoolmodel.PasswordPolicy {
	return &userpoolmodel.PasswordPolicy{
		BaseModel: sharedmodel.NewBaseModelAt(now, actorID),

		UserPoolID: userPoolID,

		MinLength: userpoolconstants.
			DefaultPasswordMinLength,

		MaxLength: userpoolconstants.
			DefaultPasswordMaxLength,

		PreventCommonPasswords: true,

		BreachedPasswordCheckEnabled: false,

		PasswordHistoryCount: userpoolconstants.
			DefaultPasswordHistoryCount,

		ResetTokenTTLSeconds: userpoolconstants.
			DefaultPasswordResetTokenTTLSeconds,

		LockoutPolicy: datatypes.JSON(
			[]byte(`{}`),
		),
	}
}

func buildDefaultRecoveryPolicy(
	userPoolID uuid.UUID,
	now time.Time,
	actorID *uuid.UUID,
) *userpoolmodel.RecoveryPolicy {
	return &userpoolmodel.RecoveryPolicy{
		BaseModel: sharedmodel.NewBaseModelAt(now, actorID),

		UserPoolID: userPoolID,

		EmailLinkEnabled: true,

		EmailOTPEnabled: true,

		SMSOTPEnabled: false,

		SavedRecoveryCodeEnabled: true,

		ExistingTOTPEnabled: true,

		ExistingPasskeyEnabled: true,

		AdminAssistedEnabled: true,

		SecurityQuestionEnabled: false,

		MinimumProofs: 1,

		RequireIndependentChannels: false,

		PostRecoverySessionAction: userpoolconstants.
			DefaultPostRecoverySessionAction,
	}
}
