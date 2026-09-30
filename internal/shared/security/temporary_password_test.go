package security

import "testing"

func TestGenerateTemporaryPassword(t *testing.T) {
	first, err := GenerateTemporaryPassword()
	if err != nil {
		t.Fatalf("GenerateTemporaryPassword() error = %v", err)
	}
	second, err := GenerateTemporaryPassword()
	if err != nil {
		t.Fatalf("GenerateTemporaryPassword() second error = %v", err)
	}
	if len(first) < 30 {
		t.Fatalf("temporary password too short: %d", len(first))
	}
	if first == second {
		t.Fatal("temporary passwords must be unique")
	}
}
