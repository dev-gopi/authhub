package constants

import (
	"testing"

	permissionconstants "github.com/dev-gopi/authhub/internal/modules/permission/constants"
)

func TestDefaultRolesAreUniqueAndContainPrimaryAdmin(t *testing.T) {
	seen := map[string]bool{}
	foundPrimary := false
	for _, role := range DefaultRoles {
		if seen[role.APILabel] {
			t.Fatalf("duplicate default role %q", role.APILabel)
		}
		seen[role.APILabel] = true
		if role.APILabel == "primary_admin" {
			foundPrimary = true
			if role.Assignable {
				t.Fatal("primary_admin must not be generally assignable")
			}
		}
	}
	if len(DefaultRoles) != 8 {
		t.Fatalf("expected 8 default roles, got %d", len(DefaultRoles))
	}
	if !foundPrimary {
		t.Fatal("primary_admin role is required")
	}
}

func TestDefaultRolePermissionsExistInCatalog(t *testing.T) {
	catalog := make(map[string]struct{}, len(permissionconstants.DefaultPermissions))
	for _, permission := range permissionconstants.DefaultPermissions {
		catalog[permission.APILabel] = struct{}{}
	}

	for _, role := range DefaultRoles {
		seen := map[string]struct{}{}
		for _, label := range role.PermissionLabels {
			if _, ok := catalog[label]; !ok {
				t.Fatalf("role %q references unknown permission %q", role.APILabel, label)
			}
			if _, duplicate := seen[label]; duplicate {
				t.Fatalf("role %q contains duplicate permission %q", role.APILabel, label)
			}
			seen[label] = struct{}{}
		}
	}
}

func TestStandardUserHasOnlySelfServicePermissions(t *testing.T) {
	for _, role := range DefaultRoles {
		if role.APILabel != "standard_user" {
			continue
		}
		for _, label := range role.PermissionLabels {
			if len(label) < len("self.") || label[:len("self.")] != "self." {
				t.Fatalf("standard_user has non-self-service permission %q", label)
			}
		}
		return
	}
	t.Fatal("standard_user role is required")
}
