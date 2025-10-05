// Package validation
package validation

import (
	"errors"
	"regexp"

	"github.com/Jrdc12/habit-tracker-go-api/internal/user"
)

var (
	errNameRequired  = errors.New("name is required")
	errEmailRequired = errors.New("email is required")
	errPassRequired  = errors.New("password is required")
	errEmailFormat   = errors.New("invalid email format")
)

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func ValidateCreateUser(r user.CreateRequest) error {
	if r.Name == "" {
		return errNameRequired
	}
	if r.Email == "" {
		return errEmailRequired
	}
	if r.Password == "" {
		return errPassRequired
	}
	if !emailRe.MatchString(r.Email) {
		return errEmailFormat
	}
	return nil
}
