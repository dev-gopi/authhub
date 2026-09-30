package service

import (
	"context"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/entity"
	"github.com/dev-gopi/authhub/internal/modules/tenant/dto"
)

type Interface interface {
	Create(
		ctx context.Context,
		req dto.CreateTenantRequest,
		rootActor *entity.RootAuthContext,
	) (*dto.CreateTenantResponse, error)
}
