package validator

import "github.com/go-playground/validator/v10"

var (
	validate *validator.Validate
)

func InitValidator() {
	validate = validator.New()
	_ = validate.RegisterValidation("role", validateRole)
}

func GetValidator() *validator.Validate {
	if validate == nil {
		InitValidator()
	}
	return validate
}

func validateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	validRoles := map[string]bool {
		"admin": true,
		"user": true,
		"mod": true,
	}
	return validRoles[role] || role == ""
}

func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}