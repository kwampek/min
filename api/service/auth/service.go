package auth

import (
	"api/models"
	"api/service/jwt"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Storage   Storage
	Jwt       *jwt.Service
	Validator *Validator
}

func NewService(Storage Storage, JwtService *jwt.Service) Service {
	return Service{
		Storage:   Storage,
		Jwt:       JwtService,
		Validator: NewValidator(),
	}
}

func (s *Service) Register(req RegisterRequest) (*AuthResponse, error) {
	if err := s.Validator.ValidateLogin(req.Login); err != nil {
		return nil, err
	}
	if err := s.Validator.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	if err := s.Validator.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	existing, _ := s.Storage.GetUserByLogin(req.Login)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Login:        req.Login,
		PasswordHash: string(hash),
		Email:        req.Email,
		CreatedAt:    time.Now(),
	}

	userID, err := s.Storage.CreateUser(user)
	if err != nil {
		return nil, err
	}

	token, err := s.Jwt.GenerateToken(userID, req.Login)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	if err := s.Storage.SaveToken(userID, token, expiresAt); err != nil {
		return nil, err
	}

	return &AuthResponse{
		UserID:      userID,
		Login:       req.Login,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Token:       token,
		CreatedAt:   time.Now(),
	}, nil
}

func (s *Service) Login(req LoginRequest) (*AuthResponse, error) {
	if err := s.Validator.ValidateLogin(req.Login); err != nil {
		return nil, err
	}
	if len(req.Password) < 1 {
		return nil, errors.New("password is required")
	}

	user, err := s.Storage.GetUserByLogin(req.Login)
	if err != nil {
		return nil, ErrUserNotFound
	}

	fmt.Println("USER: ", user)
	fmt.Println("REQ: ", req)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	token, err := s.Jwt.GenerateToken(user.UserID, user.Login)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	if err := s.Storage.SaveToken(user.UserID, token, expiresAt); err != nil {
		return nil, err
	}

	return &AuthResponse{
		UserID:      user.UserID,
		Login:       user.Login,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Token:       token,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *Service) Logout(userID int, token string) error {
	return s.Storage.DeleteToken(userID, token)
}
