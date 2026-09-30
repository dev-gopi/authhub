package service

import (
	"context"
	"time"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/constants"
	"github.com/dev-gopi/authhub/internal/modules/rootauth/entity"
	"github.com/dev-gopi/authhub/internal/shared/security"
)

func (s *Service) Authenticate(
	ctx context.Context,
	rawToken string,
) (*entity.RootAuthContext, error) {
	if rawToken == "" {
		return nil, entity.ErrUnauthorized
	}

	tokenHash := security.HashSessionToken(
		rawToken,
	)

	session, err :=
		s.sessions.FindActiveByTokenHash(
			ctx,
			tokenHash,
		)

	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, entity.ErrUnauthorized
	}

	now := time.Now().UTC()

	if now.After(session.ExpiresAt) {
		_ = s.sessions.Revoke(
			ctx,
			session.ID,
			"absolute_expired",
		)

		return nil, entity.ErrSessionExpired
	}

	if now.After(session.IdleExpiresAt) {
		_ = s.sessions.Revoke(
			ctx,
			session.ID,
			"idle_expired",
		)

		return nil, entity.ErrSessionExpired
	}

	user, err := s.users.FindByID(
		ctx,
		session.PlatformUserID,
	)
	if err != nil {
		return nil, err
	}

	if user == nil ||
		user.Status != "active" ||
		!user.IsActive ||
		user.IsDeleted {
		_ = s.sessions.Revoke(
			ctx,
			session.ID,
			"account_unavailable",
		)

		return nil, entity.ErrUnauthorized
	}

	if session.CredentialVersion !=
		user.CredentialVersion {

		_ = s.sessions.Revoke(
			ctx,
			session.ID,
			"credential_version_changed",
		)

		return nil, entity.ErrUnauthorized
	}

	if !user.IsRootAdmin {
		_ = s.sessions.Revoke(
			ctx,
			session.ID,
			"root_access_removed",
		)

		return nil, entity.ErrForbidden
	}

	if err := s.sessions.Touch(
		ctx,
		session.ID,
		now,
		now.Add(
			constants.SessionIdleLifetime,
		),
	); err != nil {
		return nil, err
	}

	return &entity.RootAuthContext{
		PlatformUserID: user.ID,
		SessionID:      session.ID,

		CredentialVersion: user.CredentialVersion,

		IsRootAdmin: user.IsRootAdmin,
	}, nil
}
