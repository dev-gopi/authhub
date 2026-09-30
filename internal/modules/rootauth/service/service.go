package bootstrap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	platformmodel "github.com/dev-gopi/authhub/internal/modules/platformuser/model"
	platformrepository "github.com/dev-gopi/authhub/internal/modules/platformuser/repository"
	rootmodel "github.com/dev-gopi/authhub/internal/modules/rootauth/model"
	rootrepository "github.com/dev-gopi/authhub/internal/modules/rootauth/repository"
	"github.com/dev-gopi/authhub/internal/platform/database"
	"github.com/dev-gopi/authhub/internal/shared/security"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db *database.Postgres

	platformUsers platformrepository.Repository
	passwords     rootrepository.PasswordRepository
	hasher        *security.PasswordHasher
}

func NewService(
	db *database.Postgres,
	platformUsers platformrepository.Repository,
	passwords rootrepository.PasswordRepository,
	hasher *security.PasswordHasher,
) *Service {
	return &Service{
		db:            db,
		platformUsers: platformUsers,
		passwords:     passwords,
		hasher:        hasher,
	}
}

func (s *Service) CreateRootAdmin(
	ctx context.Context,
	req CreateRootAdminRequest,
) (*platformmodel.PlatformUser, error) {
	username := strings.TrimSpace(req.Username)
	email := strings.ToLower(
		strings.TrimSpace(req.Email),
	)
	displayName := strings.TrimSpace(req.DisplayName)

	if username == "" {
		return nil, errors.New(
			"username is required",
		)
	}

	if email == "" {
		return nil, errors.New(
			"email is required",
		)
	}

	if displayName == "" {
		displayName = username
	}

	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	exists, err := s.platformUsers.ExistsByUsernameOrEmail(
		ctx,
		username,
		email,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"check existing platform user: %w",
			err,
		)
	}

	if exists {
		return nil, errors.New(
			"platform user with username or email already exists",
		)
	}

	passwordHash, err := s.hasher.Hash(
		req.Password,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"hash root admin password: %w",
			err,
		)
	}

	paramsJSON, err := json.Marshal(
		s.hasher.Params(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"marshal password parameters: %w",
			err,
		)
	}

	now := time.Now().UTC()

	user := &platformmodel.PlatformUser{
		ID: uuid.New(),

		Username: username,
		Email:    email,

		DisplayName: displayName,

		Status: "active",

		IsRootAdmin: true,

		CredentialVersion: 1,

		CreatedAt: now,
		UpdatedAt: now,

		IsActive:  true,
		IsDeleted: false,
	}

	password := &rootmodel.PlatformPassword{
		ID: uuid.New(),

		PlatformUserID: user.ID,

		PasswordHash: passwordHash,

		PasswordAlgorithm: "argon2id",

		PasswordParams: paramsJSON,

		PasswordVersion: 1,

		ChangedAt: now,

		CreatedAt: now,
		UpdatedAt: now,

		IsActive:  true,
		IsDeleted: false,
	}

	err = s.db.Transaction(
		ctx,
		func(tx *gorm.DB) error {
			if err := s.platformUsers.Create(
				ctx,
				tx,
				user,
			); err != nil {
				return fmt.Errorf(
					"create root platform user: %w",
					err,
				)
			}

			if err := s.passwords.Create(
				ctx,
				tx,
				password,
			); err != nil {
				return fmt.Errorf(
					"create root password: %w",
					err,
				)
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
