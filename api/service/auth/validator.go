package auth

import (
	"errors"
	"regexp"
)

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateLogin(login string) error {
	return nil
}

func (v *Validator) ValidatePassword(password string) error {
	return nil
}

func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return nil
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if !matched {
		return errors.New("invalid email format")
	}

	return nil
}
