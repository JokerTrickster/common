package validator

/*
	커스텀 검증 규칙 정의
*/

import (
	"github.com/go-playground/validator/v10"
)

// RegisterCustomRules registers custom validation rules
func RegisterCustomRules(v *validator.Validate) {
	// Example: Register a custom validation rule
	_ = v.RegisterValidation("customRule", func(fl validator.FieldLevel) bool {
		// Add your custom rule logic here
		return len(fl.Field().String()) > 5 // Example: Field must be longer than 5 characters
	})
}

// InitValidator initializes the validator instance with custom rules
func InitValidator() *validator.Validate {
	v := validator.New()
	RegisterCustomRules(v)
	return v
}
