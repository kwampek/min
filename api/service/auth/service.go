package auth

import (
	"api/models"
	"api/service/tokens"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Storage      Storage
	TokenService tokens.TokenService
	Validator    Validator
}

func NewAuthService(Storage Storage, TokenService tokens.TokenService) AuthService {
	return AuthService{
		Storage:      Storage,
		TokenService: TokenService,
		Validator:    NewValidator(),
	}
}

func GetSessionInfo(r *http.Request, userID int, token string) (models.Session, error) {
	hash := sha256.Sum256([]byte(token))

	return models.Session{
		UserID:     userID,
		TokenHash:  hex.EncodeToString(hash[:]),
		DeviceName: r.UserAgent(),
		IPAddress:  r.RemoteAddr,
		ExpiresAt:  time.Now().Add(30 * 24 * time.Hour),
	}, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) Register(req RegisterRequest, device models.Device) (*AuthResponse, error) {
	if err := s.Validator.ValidateLogin(req.Login); err != nil {
		return nil, err
	}
	if err := s.Validator.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	if err := s.Validator.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	_, err := s.Storage.GetUserByLogin(req.Login)
	if err != sql.ErrNoRows {
		return nil, ErrUserAlreadyExists
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Login:        req.Login,
		PasswordHash: string(hashPassword),
		PhoneNumber:  req.PhoneNumber,
		Email:        req.Email,
		CreatedAt:    time.Now(),
	}

	userID, err := s.Storage.CreateUser(user)
	if err != nil {
		return nil, err
	}

	token, err := s.TokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	err = s.Storage.CreateSession(models.Session{
		UserID:     userID,
		TokenHash:  hashToken(token),
		DeviceName: device.DeviceName,
		IPAddress:  device.IPAddress,
		ExpiresAt:  time.Now().Add(30 * 24 * time.Hour),
	})
	if err != nil {
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

func (s *AuthService) Login(req LoginRequest, device models.Device) (*AuthResponse, error) {
	if err := s.Validator.ValidateLogin(req.Login); err != nil {
		return nil, err
	}

	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	user, err := s.Storage.GetUserByLogin(req.Login)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, ErrInvalidPassword
	}

	token, err := s.TokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	err = s.Storage.CreateSession(models.Session{
		UserID:     user.UserID,
		TokenHash:  hashToken(token),
		DeviceName: device.DeviceName,
		IPAddress:  device.IPAddress,
		ExpiresAt:  time.Now().Add(30 * 24 * time.Hour),
	})
	if err != nil {
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

func (s *AuthService) Logout(token string) error {
	return s.Storage.RevokeSession(hashToken(token))
}
