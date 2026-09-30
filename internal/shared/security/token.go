package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

const sessionTokenBytes = 32

func GenerateSessionToken() (string, string, error) {
	randomBytes := make([]byte, sessionTokenBytes)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf(
			"generate secure session token: %w",
			err,
		)
	}

	rawToken := base64.RawURLEncoding.EncodeToString(
		randomBytes,
	)

	tokenHash := HashSessionToken(rawToken)

	return rawToken, tokenHash, nil
}

func HashSessionToken(token string) string {
	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}

func HashIdentifier(value string) string {
	normalized := strings.ToLower(
		strings.TrimSpace(value),
	)

	hash := sha256.Sum256(
		[]byte(normalized),
	)

	return hex.EncodeToString(hash[:])
}
