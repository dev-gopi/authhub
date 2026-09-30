package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dev-gopi/authhub/internal/shared/security"
	"github.com/google/uuid"
)

type preparedTemporaryCredential struct {
	PasswordHash       string
	PasswordParams     []byte
	EncryptedPlaintext string
}

func (s *Service) prepareTemporaryCredential(ctx context.Context, tenantID uuid.UUID) (*preparedTemporaryCredential, error) {
	plain, err := security.GenerateTemporaryPassword()
	if err != nil {
		return nil, err
	}
	defer func() { plain = "" }()

	hash, err := s.passwordHasher.Hash(plain)
	if err != nil {
		return nil, fmt.Errorf("hash temporary password: %w", err)
	}
	params, err := json.Marshal(s.passwordHasher.Params())
	if err != nil {
		return nil, fmt.Errorf("marshal password params: %w", err)
	}
	ciphertext, err := s.keyManagement.EncryptTenantData(ctx, tenantID, []byte(plain))
	if err != nil {
		return nil, fmt.Errorf("encrypt temporary password for delivery: %w", err)
	}
	return &preparedTemporaryCredential{PasswordHash: hash, PasswordParams: params, EncryptedPlaintext: ciphertext}, nil
}
