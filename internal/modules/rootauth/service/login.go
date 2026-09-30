package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/constants"
	"github.com/dev-gopi/authhub/internal/modules/rootauth/dto"
	"github.com/dev-gopi/authhub/internal/modules/rootauth/entity"
	rootmodel "github.com/dev-gopi/authhub/internal/modules/rootauth/model"
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"github.com/dev-gopi/authhub/internal/shared/security"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) Login(
	ctx context.Context,
	req dto.LoginRequest,
	metadata entity.RequestMetadata,
) (*dto.LoginResponse, error) {
	if err := s.checkRateLimit(
		ctx,
		req.Identifier,
		metadata.IPAddress,
	); err != nil {
		s.recordLoginAttempt(
			ctx,
			nil,
			req.Identifier,
			metadata,
			false,
			"rate_limited",
		)

		return nil, err
	}

	user, err := s.users.FindByIdentifier(
		ctx,
		req.Identifier,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"find platform user: %w",
			err,
		)
	}

	if user == nil {
		_, _ = s.passwordHasher.Verify(
			req.Password,
			s.dummyPasswordHash,
		)

		s.recordRateLimitFailure(
			ctx,
			req.Identifier,
			metadata.IPAddress,
		)

		s.recordLoginAttempt(
			ctx,
			nil,
			req.Identifier,
			metadata,
			false,
			"invalid_credentials",
		)

		return nil, entity.ErrInvalidCredentials
	}

	userID := user.ID

	if user.Status != "active" ||
		!user.IsActive ||
		user.IsDeleted {
		s.recordRateLimitFailure(
			ctx,
			req.Identifier,
			metadata.IPAddress,
		)

		s.recordLoginAttempt(
			ctx,
			&userID,
			req.Identifier,
			metadata,
			false,
			"account_disabled",
		)

		return nil, entity.ErrInvalidCredentials
	}

	if !user.IsRootAdmin {
		s.recordRateLimitFailure(
			ctx,
			req.Identifier,
			metadata.IPAddress,
		)

		s.recordLoginAttempt(
			ctx,
			&userID,
			req.Identifier,
			metadata,
			false,
			"not_root_admin",
		)

		return nil, entity.ErrInvalidCredentials
	}

	password, err :=
		s.passwords.FindByPlatformUserID(
			ctx,
			user.ID,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"find root password: %w",
			err,
		)
	}

	if password == nil {
		_, _ = s.passwordHasher.Verify(
			req.Password,
			s.dummyPasswordHash,
		)

		s.recordRateLimitFailure(
			ctx,
			req.Identifier,
			metadata.IPAddress,
		)

		s.recordLoginAttempt(
			ctx,
			&userID,
			req.Identifier,
			metadata,
			false,
			"invalid_credentials",
		)

		return nil, entity.ErrInvalidCredentials
	}

	valid, err := s.passwordHasher.Verify(
		req.Password,
		password.PasswordHash,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"verify root password: %w",
			err,
		)
	}

	if !valid {
		s.recordRateLimitFailure(
			ctx,
			req.Identifier,
			metadata.IPAddress,
		)

		s.recordLoginAttempt(
			ctx,
			&userID,
			req.Identifier,
			metadata,
			false,
			"invalid_credentials",
		)

		return nil, entity.ErrInvalidCredentials
	}

	needsRehash, err :=
		s.passwordHasher.NeedsRehash(
			password.PasswordHash,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"check root password rehash: %w",
			err,
		)
	}

	if needsRehash {
		newHash, err :=
			s.passwordHasher.Hash(
				req.Password,
			)

		if err != nil {
			return nil, fmt.Errorf(
				"rehash root password: %w",
				err,
			)
		}

		params, err := json.Marshal(
			s.passwordHasher.Params(),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"marshal password params: %w",
				err,
			)
		}

		if err := s.passwords.UpdateHash(
			ctx,
			s.db,
			user.ID,
			newHash,
			params,
		); err != nil {
			return nil, fmt.Errorf(
				"update root password hash: %w",
				err,
			)
		}
	}

	rawToken, tokenHash, err :=
		security.GenerateSessionToken()

	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	var ip *string
	if metadata.IPAddress != "" {
		ip = &metadata.IPAddress
	}

	var userAgent *string
	if metadata.UserAgent != "" {
		userAgent = &metadata.UserAgent
	}

	session := &rootmodel.PlatformSession{
		BaseModel: sharedmodel.BaseModel{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			IsActive:  true,
			IsDeleted: false,
		},

		PlatformUserID: user.ID,

		SessionTokenHash: tokenHash,

		CredentialVersion: user.CredentialVersion,

		IPAddress: ip,

		UserAgent: userAgent,

		ExpiresAt: now.Add(
			constants.SessionAbsoluteLifetime,
		),

		IdleExpiresAt: now.Add(
			constants.SessionIdleLifetime,
		),

		LastSeenAt: now,
	}

	if err := s.sessions.Create(
		ctx,
		session,
	); err != nil {
		return nil, fmt.Errorf(
			"create root session: %w",
			err,
		)
	}

	if err := s.users.UpdateLastLogin(
		ctx,
		user.ID,
		now,
	); err != nil {
		_ = s.sessions.Revoke(
			ctx,
			session.ID,
			"login_finalization_failed",
		)

		return nil, fmt.Errorf(
			"update root last login: %w",
			err,
		)
	}

	s.clearRateLimit(
		ctx,
		req.Identifier,
		metadata.IPAddress,
	)

	s.recordLoginAttempt(
		ctx,
		&userID,
		req.Identifier,
		metadata,
		true,
		"",
	)

	s.logger.Info(
		"root login succeeded",
		zap.String(
			"platform_user_id",
			user.ID.String(),
		),
		zap.String(
			"session_id",
			session.ID.String(),
		),
	)

	return &dto.LoginResponse{
		SessionToken: rawToken,
		TokenType:    constants.SessionTokenType,
		ExpiresAt:    session.ExpiresAt,
	}, nil
}

func (s *Service) recordLoginAttempt(
	ctx context.Context,
	userID *uuid.UUID,
	identifier string,
	metadata entity.RequestMetadata,
	success bool,
	failureReason string,
) {
	now := time.Now().UTC()

	var reason *string
	if failureReason != "" {
		reason = &failureReason
	}

	var ip *string
	if metadata.IPAddress != "" {
		ip = &metadata.IPAddress
	}

	var userAgent *string
	if metadata.UserAgent != "" {
		userAgent = &metadata.UserAgent
	}

	attempt := &rootmodel.PlatformLoginAttempt{
		BaseModel: sharedmodel.BaseModel{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			IsActive:  true,
			IsDeleted: false,
		},

		PlatformUserID: userID,

		IdentifierHash: security.HashIdentifier(
			identifier,
		),

		IPAddress: ip,

		UserAgent: userAgent,

		Success: success,

		FailureReason: reason,

		AttemptedAt: now,
	}

	if err := s.loginAttempts.Create(
		ctx,
		attempt,
	); err != nil {
		s.logger.Error(
			"failed to persist root login attempt",
			zap.Error(err),
		)
	}
}
