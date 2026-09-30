package validator

import (
	"errors"
	"strings"
	"unicode"
)

func NormalizeIdentifier(identifier string) string {
	return strings.ToLower(
		strings.TrimSpace(identifier),
	)
}

func validatePassword(password string) error {
	if len(password) < 14 {
		return errors.New(
			"root admin password must contain at least 14 characters",
		)
	}

	if len(password) > 256 {
		return errors.New(
			"root admin password exceeds maximum length",
		)
	}

	var hasUpper bool
	var hasLower bool
	var hasNumber bool

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true

		case unicode.IsLower(ch):
			hasLower = true

		case unicode.IsNumber(ch):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return errors.New(
			"root admin password must contain uppercase, lowercase and numeric characters",
		)
	}

	if strings.Contains(
		strings.ToLower(password),
		"password",
	) {
		return errors.New(
			"root admin password is too weak",
		)
	}

	return nil
}
