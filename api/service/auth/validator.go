package auth

import (
	"api/models"
	"errors"
	"regexp"
)

type Validator struct {
}

func NewValidator() Validator {
	return Validator{}
}

func (v *Validator) ValidateLogin(login string) error {
	if login == "" {
		return errors.New("login is required")
	}
	if len(login) > 32 {
		return errors.New("login must be at most 32 characters")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._-]+$`, login)
	if !matched {
		return errors.New("login contains invalid characters")
	}
	return nil
}

func (v *Validator) ValidatePassword(password string) error {
	// CLEAR IT

	// if password == "" {
	// 	return errors.New("password is required")
	// }
	// if len(password) < 8 {
	// 	return errors.New("password must be at least 8 characters")
	// }
	// if len(password) > 72 {
	// 	return errors.New("password must be at most 72 characters")
	// }

	return nil
}

func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return nil
	}
	if len(email) > 254 {
		return errors.New("email must be at most 254 characters")
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func (v *Validator) ValidatePhone(phone string) error {
	if phone == "" {
		return nil
	}
	if len(phone) > 20 {
		return errors.New("phone number must be at most 20 characters")
	}

	matched, _ := regexp.MatchString(`^\+?[0-9]+$`, phone)
	if !matched {
		return errors.New("phone number contains invalid characters")
	}
	return nil
}

func (v *Validator) ValidateRegister(req RegisterRequest) error {
	if err := v.ValidateLogin(req.Login); err != nil {
		return err
	}
	if err := v.ValidatePassword(req.Password); err != nil {
		return err
	}
	if err := v.ValidateEmail(req.Email); err != nil {
		return err
	}
	if err := v.ValidatePhone(req.PhoneNumber); err != nil {
		return err
	}
	return nil
}

func (v *Validator) Validate(req LoginRequest) error {
	if err := v.ValidateLogin(req.Login); err != nil {
		return err
	}
	if err := v.ValidatePassword(req.Password); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ValidateSession(token string) (models.Identifier, error) {
	return s.Storage.ValidateSession(token)
}
