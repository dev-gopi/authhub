package validator

import (
	"errors"

	govalidator "github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *govalidator.Validate
}

func New() *Validator {
	return &Validator{
		validate: govalidator.New(),
	}
}

func (v *Validator) Struct(value any) error {
	if err := v.validate.Struct(value); err != nil {
		var validationErrors govalidator.ValidationErrors

		if errors.As(err, &validationErrors) {
			return validationErrors
		}

		return err
	}

	return nil
}
