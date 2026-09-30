package service

import (
	"context"
	vaultinfra "github.com/dev-gopi/authhub/internal/platform/vault"
)

type TransitClient interface {
	EnsureTransitKey(ctx context.Context, mount string, spec vaultinfra.TransitKeySpec) error
	EncryptTransit(ctx context.Context, mount string, keyName string, plaintext []byte) (string, error)
}
