package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatErrors(err error) string {
	var errors string

	for _, err := range err.(validator.ValidationErrors) {
		errors += fmt.Sprintf("Field '%s' failed validation on '%s'; ",
			err.Field(), err.Tag())
	}

	return errors
}
