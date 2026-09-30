package service

import (
	"context"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/dto"
	"github.com/dev-gopi/authhub/internal/modules/rootauth/entity"
)

func (s *Service) Me(
	ctx context.Context,
	auth *entity.RootAuthContext,
) (*dto.MeResponse, error) {
	user, err := s.users.FindByID(
		ctx,
		auth.PlatformUserID,
	)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, entity.ErrUnauthorized
	}

	return &dto.MeResponse{
		ID: user.ID,

		Username: user.Username,

		Email: user.Email,

		DisplayName: user.DisplayName,

		IsRootAdmin: user.IsRootAdmin,
	}, nil
}
