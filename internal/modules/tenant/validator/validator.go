package validator

import (
	"errors"
	"regexp"
	"strings"

	"github.com/dev-gopi/authhub/internal/modules/tenant/dto"
	sharedvalidator "github.com/dev-gopi/authhub/internal/shared/validator"
)

var apiLabelPattern = regexp.MustCompile(
	`^[a-z][a-z0-9_]{1,62}$`,
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

func (v *Validator) ValidateCreate(
	req *dto.CreateTenantRequest,
) error {
	req.APILabel = strings.ToLower(
		strings.TrimSpace(req.APILabel),
	)

	req.DisplayName = strings.TrimSpace(
		req.DisplayName,
	)

	req.PrimaryAdmin.Username = strings.ToLower(
		strings.TrimSpace(
			req.PrimaryAdmin.Username,
		),
	)

	req.PrimaryAdmin.Email = strings.ToLower(
		strings.TrimSpace(
			req.PrimaryAdmin.Email,
		),
	)

	req.PrimaryAdmin.DisplayName =
		strings.TrimSpace(
			req.PrimaryAdmin.DisplayName,
		)

	if err := v.validator.Struct(req); err != nil {
		return err
	}

	if !apiLabelPattern.MatchString(
		req.APILabel,
	) {
		return errors.New(
			"invalid tenant api_label",
		)
	}

	return nil
}
