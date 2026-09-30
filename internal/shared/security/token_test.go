package security

import "testing"

func TestGenerateSessionToken(t *testing.T) {
	rawToken, tokenHash, err :=
		GenerateSessionToken()

	if err != nil {
		t.Fatalf(
			"generate session token: %v",
			err,
		)
	}

	if rawToken == "" {
		t.Fatal("raw token is empty")
	}

	if tokenHash == "" {
		t.Fatal("token hash is empty")
	}

	if rawToken == tokenHash {
		t.Fatal(
			"raw token must not equal token hash",
		)
	}

	if HashSessionToken(rawToken) != tokenHash {
		t.Fatal(
			"session token hash is not deterministic",
		)
	}
}
