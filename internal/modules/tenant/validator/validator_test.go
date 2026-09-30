package validator

import (
	"testing"

	"github.com/dev-gopi/authhub/internal/modules/tenant/dto"
	sharedvalidator "github.com/dev-gopi/authhub/internal/shared/validator"
)

func TestValidateCreateNormalizesTenantAndAdminIdentifiers(t *testing.T) {
	v := New(sharedvalidator.New())
	req := &dto.CreateTenantRequest{
		APILabel:    "Acme",
		DisplayName: " Acme Corporation ",
		Profile:     dto.TenantProfileRequest{},
		PrimaryAdmin: dto.PrimaryAdminRequest{
			Username:    " ACME_ADMIN ",
			Email:       " ADMIN@EXAMPLE.COM ",
			DisplayName: " Acme Admin ",
		},
	}
	if err := v.ValidateCreate(req); err != nil {
		t.Fatalf("ValidateCreate() error = %v", err)
	}
	if req.APILabel != "acme" {
		t.Fatalf("api label not normalized: %q", req.APILabel)
	}
	if req.PrimaryAdmin.Username != "acme_admin" {
		t.Fatalf("username not normalized: %q", req.PrimaryAdmin.Username)
	}
	if req.PrimaryAdmin.Email != "admin@example.com" {
		t.Fatalf("email not normalized: %q", req.PrimaryAdmin.Email)
	}
}

func TestValidateCreateRejectsUnsafeAPILabel(t *testing.T) {
	v := New(sharedvalidator.New())
	req := &dto.CreateTenantRequest{
		APILabel:     "bad-label!",
		DisplayName:  "Bad Tenant",
		PrimaryAdmin: dto.PrimaryAdminRequest{Username: "admin", Email: "admin@example.com", DisplayName: "Admin"},
	}
	if err := v.ValidateCreate(req); err == nil {
		t.Fatal("expected invalid api_label error")
	}
}
