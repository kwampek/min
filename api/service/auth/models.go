package auth

import (
	"errors"
	"time"
)

type Device struct {
	DeviceName string
	IPAddress  string
}

type RegisterRequest struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	UserID      int       `json:"user_id"`
	Login       string    `json:"login"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Token       string    `json:"token"`
	CreatedAt   time.Time `json:"created_at"`
}

type TokenInfo struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

var (
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrWeakPassword      = errors.New("password is too weak")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrLoginTooShort     = errors.New("login must be at least 3 characters")
	ErrLoginTooLong      = errors.New("login must be at most 32 characters")
)
