package service

import (
	auditrepository "github.com/dev-gopi/authhub/internal/modules/audit/repository"
	keymanagement "github.com/dev-gopi/authhub/internal/modules/keymanagement/service"
	outboxrepository "github.com/dev-gopi/authhub/internal/modules/outbox/repository"
	permissionrepository "github.com/dev-gopi/authhub/internal/modules/permission/repository"
	platformuserrepository "github.com/dev-gopi/authhub/internal/modules/platformuser/repository"
	rolerepository "github.com/dev-gopi/authhub/internal/modules/role/repository"
	rootrepository "github.com/dev-gopi/authhub/internal/modules/rootauth/repository"
	tenantrepository "github.com/dev-gopi/authhub/internal/modules/tenant/repository"
	tenantmemberrepository "github.com/dev-gopi/authhub/internal/modules/tenantmember/repository"
	userpoolrepository "github.com/dev-gopi/authhub/internal/modules/userpool/repository"
	"github.com/dev-gopi/authhub/internal/shared/security"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB

	tenants       tenantrepository.Interface
	profiles      tenantrepository.ProfileRepository
	userPools     userpoolrepository.Interface
	permissions   permissionrepository.Interface
	roles         rolerepository.Interface
	platformUsers platformuserrepository.Interface
	tenantMembers tenantmemberrepository.Interface
	passwords     rootrepository.PasswordRepository
	outbox        outboxrepository.Interface
	audit         auditrepository.Interface

	passwordHasher             *security.PasswordHasher
	keyManagement              keymanagement.Interface
	issuerBaseURL              string
	defaultRoleTemplateVersion string
	logger                     *zap.Logger
}

func New(
	db *gorm.DB,
	tenants tenantrepository.Interface,
	profiles tenantrepository.ProfileRepository,
	userPools userpoolrepository.Interface,
	permissions permissionrepository.Interface,
	roles rolerepository.Interface,
	platformUsers platformuserrepository.Interface,
	tenantMembers tenantmemberrepository.Interface,
	passwords rootrepository.PasswordRepository,
	outbox outboxrepository.Interface,
	audit auditrepository.Interface,
	passwordHasher *security.PasswordHasher,
	keyManagement keymanagement.Interface,
	issuerBaseURL string,
	defaultRoleTemplateVersion string,
	logger *zap.Logger,
) *Service {
	return &Service{
		db:                         db,
		tenants:                    tenants,
		profiles:                   profiles,
		userPools:                  userPools,
		permissions:                permissions,
		roles:                      roles,
		platformUsers:              platformUsers,
		tenantMembers:              tenantMembers,
		passwords:                  passwords,
		outbox:                     outbox,
		audit:                      audit,
		passwordHasher:             passwordHasher,
		keyManagement:              keyManagement,
		issuerBaseURL:              issuerBaseURL,
		defaultRoleTemplateVersion: defaultRoleTemplateVersion,
		logger:                     logger,
	}
}
