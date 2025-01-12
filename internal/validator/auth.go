package validator

import (
	"fmt"
	"pulse/graph/graphql_model"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	validate   *validator.Validate
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func init() {
	validate = validator.New()
}

func ValidateRegisterInput(input graphql_model.RegisterInput) error {
	if !emailRegex.MatchString(input.Email) {
		return fmt.Errorf("invalid email format")
	}

	if len(input.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	if len(input.Name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}

	return nil
}
