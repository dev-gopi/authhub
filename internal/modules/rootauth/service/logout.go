package service

import (
	"context"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/entity"
)

func (s *Service) Logout(
	ctx context.Context,
	auth *entity.RootAuthContext,
) error {
	return s.sessions.Revoke(
		ctx,
		auth.SessionID,
		"user_logout",
	)
}

func (s *Service) LogoutAll(
	ctx context.Context,
	auth *entity.RootAuthContext,
) error {
	return s.sessions.RevokeAllByPlatformUser(
		ctx,
		auth.PlatformUserID,
		"user_logout_all",
	)
}
