package service

import (
	"fmt"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/repository"
	redisinfra "github.com/dev-gopi/authhub/internal/platform/redis"
	"github.com/dev-gopi/authhub/internal/shared/security"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB

	redis *redisinfra.Client

	users repository.PlatformUserRepository

	passwords repository.PasswordRepository

	sessions repository.SessionRepository

	loginAttempts repository.LoginAttemptRepository

	passwordHasher *security.PasswordHasher

	dummyPasswordHash string

	logger *zap.Logger
}

func New(
	db *gorm.DB,
	redis *redisinfra.Client,
	users repository.PlatformUserRepository,
	passwords repository.PasswordRepository,
	sessions repository.SessionRepository,
	loginAttempts repository.LoginAttemptRepository,
	passwordHasher *security.PasswordHasher,
	logger *zap.Logger,
) (*Service, error) {
	dummyPasswordHash, err := passwordHasher.Hash(
		"AuthHubDummyPassword123456789!",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create dummy password hash: %w",
			err,
		)
	}

	return &Service{
		db:                db,
		redis:             redis,
		users:             users,
		passwords:         passwords,
		sessions:          sessions,
		loginAttempts:     loginAttempts,
		passwordHasher:    passwordHasher,
		dummyPasswordHash: dummyPasswordHash,
		logger:            logger,
	}, nil
}
