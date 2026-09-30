package security

import "testing"

func TestPasswordHasher(t *testing.T) {
	hasher := NewPasswordHasher()

	password := "VeryStrongTestPassword!123"

	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	if hash == password {
		t.Fatal("password must never equal stored hash")
	}

	ok, err := hasher.Verify(password, hash)
	if err != nil {
		t.Fatalf("verify password: %v", err)
	}

	if !ok {
		t.Fatal("valid password was rejected")
	}

	ok, err = hasher.Verify(
		"WrongPassword",
		hash,
	)
	if err != nil {
		t.Fatalf("verify wrong password: %v", err)
	}

	if ok {
		t.Fatal("wrong password was accepted")
	}

	needsRehash, err := hasher.NeedsRehash(hash)
	if err != nil {
		t.Fatalf("check rehash: %v", err)
	}

	if needsRehash {
		t.Fatal("new hash unexpectedly requires rehash")
	}
}
