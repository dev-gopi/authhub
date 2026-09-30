package service

import (
	"context"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/dto"
	"github.com/dev-gopi/authhub/internal/modules/rootauth/entity"
)

type Interface interface {
	Login(
		ctx context.Context,
		req dto.LoginRequest,
		metadata entity.RequestMetadata,
	) (*dto.LoginResponse, error)

	Authenticate(
		ctx context.Context,
		rawToken string,
	) (*entity.RootAuthContext, error)

	Logout(
		ctx context.Context,
		auth *entity.RootAuthContext,
	) error

	LogoutAll(
		ctx context.Context,
		auth *entity.RootAuthContext,
	) error

	Me(
		ctx context.Context,
		auth *entity.RootAuthContext,
	) (*dto.MeResponse, error)
}
