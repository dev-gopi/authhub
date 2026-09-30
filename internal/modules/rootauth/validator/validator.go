package validator

import (
	"strings"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/dto"
	sharedvalidator "github.com/dev-gopi/authhub/internal/shared/validator"
)

type Validator struct {
	validator *sharedvalidator.Validator
}

func New(
	validator *sharedvalidator.Validator,
) *Validator {
	return &Validator{
		validator: validator,
	}
}

func NormalizeIdentifier(
	identifier string,
) string {
	return strings.ToLower(
		strings.TrimSpace(identifier),
	)
}

func (v *Validator) ValidateLogin(
	req *dto.LoginRequest,
) error {
	req.Identifier = NormalizeIdentifier(
		req.Identifier,
	)

	return v.validator.Struct(req)
}
