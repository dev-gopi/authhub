package constants

import "testing"

func TestDefaultPermissionsAreUniqueAndComplete(t *testing.T) {
	seen := make(map[string]struct{}, len(DefaultPermissions))
	for _, permission := range DefaultPermissions {
		if permission.APILabel == "" || permission.DisplayName == "" || permission.Resource == "" || permission.Action == "" {
			t.Fatalf("invalid default permission: %+v", permission)
		}
		if _, exists := seen[permission.APILabel]; exists {
			t.Fatalf("duplicate default permission %q", permission.APILabel)
		}
		seen[permission.APILabel] = struct{}{}
	}
	if len(DefaultPermissions) < 60 {
		t.Fatalf("expected detailed permission catalog, got only %d permissions", len(DefaultPermissions))
	}
}
