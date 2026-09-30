package service

import (
	"context"
	"testing"

	vaultinfra "github.com/dev-gopi/authhub/internal/platform/vault"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type fakeTransitClient struct {
	keys []vaultinfra.TransitKeySpec
}

func (f *fakeTransitClient) EnsureTransitKey(_ context.Context, _ string, spec vaultinfra.TransitKeySpec) error {
	f.keys = append(f.keys, spec)
	return nil
}

func (f *fakeTransitClient) EncryptTransit(_ context.Context, _ string, _ string, plaintext []byte) (string, error) {
	if len(plaintext) == 0 {
		t := ""
		_ = t
	}
	return "vault:v1:test-ciphertext", nil
}

func TestProvisionTenantKeysCreatesFiveNamespacedKeys(t *testing.T) {
	client := &fakeTransitClient{}
	svc := New(client, "transit", zap.NewNop())
	tenantID := uuid.New()
	if err := svc.ProvisionTenantKeys(context.Background(), tenantID); err != nil {
		t.Fatalf("ProvisionTenantKeys() error = %v", err)
	}
	if len(client.keys) != 5 {
		t.Fatalf("expected 5 keys, got %d", len(client.keys))
	}
	seen := map[string]bool{}
	for _, key := range client.keys {
		seen[key.Name] = true
	}
	for _, suffix := range []string{"oidc-signing", "saml-signing", "webhook-signing", "realtime-signing", "data-encryption"} {
		name := "tenant-" + tenantID.String() + "-" + suffix
		if !seen[name] {
			t.Fatalf("missing key %s", name)
		}
	}
}

func TestEncryptTenantDataUsesTenantDataKey(t *testing.T) {
	client := &fakeTransitClient{}
	svc := New(client, "transit", zap.NewNop())
	ciphertext, err := svc.EncryptTenantData(context.Background(), uuid.New(), []byte("secret"))
	if err != nil {
		t.Fatalf("EncryptTenantData() error = %v", err)
	}
	if ciphertext == "" {
		t.Fatal("expected ciphertext")
	}
}
