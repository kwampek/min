package auth

import (
	"api/models"
	"api/service/jwt"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	storage   Storage
	jwt       *jwt.Service
	validator *Validator
}

func NewService(storage Storage, jwtService *jwt.Service) *Service {
	return &Service{
		storage:   storage,
		jwt:       jwtService,
		validator: NewValidator(),
	}
}

func (s *Service) Register(req RegisterRequest) (*AuthResponse, error) {
	if err := s.validator.ValidateLogin(req.Login); err != nil {
		return nil, err
	}
	if err := s.validator.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	if err := s.validator.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	existing, _ := s.storage.GetUserByLogin(req.Login)
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

	userID, err := s.storage.CreateUser(user)
	if err != nil {
		return nil, err
	}

	token, err := s.jwt.GenerateToken(userID, req.Login)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	if err := s.storage.SaveToken(userID, token, expiresAt); err != nil {
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
	if err := s.validator.ValidateLogin(req.Login); err != nil {
		return nil, err
	}
	if len(req.Password) < 1 {
		return nil, errors.New("password is required")
	}

	user, err := s.storage.GetUserByLogin(req.Login)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwt.GenerateToken(user.UserID, user.Login)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	if err := s.storage.SaveToken(user.UserID, token, expiresAt); err != nil {
		return nil, err
	}

	return &AuthResponse{
		UserID:      user.UserID,
		Login:       user.Login,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Token:       token,
	}, nil
}

func (s *Service) Logout(userID int, token string) error {
	return s.storage.DeleteToken(userID, token)
}
