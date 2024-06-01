package validators

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates any struct passed to it.
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ValidateVariable validates a single variable with the provided tag.
func ValidateVariable(variable interface{}, tag string) error {
	return validate.Var(variable, tag)
}
