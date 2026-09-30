package service

import (
	keymanagement "github.com/dev-gopi/authhub/internal/modules/keymanagement/service"
	tenantrepository "github.com/dev-gopi/authhub/internal/modules/tenant/repository"
	userpoolrepository "github.com/dev-gopi/authhub/internal/modules/userpool/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB

	tenants tenantrepository.Interface

	profiles tenantrepository.ProfileRepository

	userPools userpoolrepository.Interface

	keyManagement keymanagement.Interface

	issuerBaseURL string

	logger *zap.Logger
}

func New(
	db *gorm.DB,
	tenants tenantrepository.Interface,
	profiles tenantrepository.ProfileRepository,
	userPools userpoolrepository.Interface,
	keyManagement keymanagement.Interface,
	issuerBaseURL string,
	logger *zap.Logger,
) *Service {
	return &Service{
		db: db,

		tenants: tenants,

		profiles: profiles,

		userPools: userPools,

		keyManagement: keyManagement,

		issuerBaseURL: issuerBaseURL,

		logger: logger,
	}
}
