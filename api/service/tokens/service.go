package tokens

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret     []byte
	issuer     string
	expiryDays int
}

type Claims struct {
	UserID int    `json:"user_id"`
	Login  string `json:"login"`
	Email  string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

func NewTokenService(secret string, issuer string, expiryDays int) TokenService {
	return TokenService{
		secret:     []byte(secret),
		issuer:     issuer,
		expiryDays: expiryDays,
	}
}

func (s *TokenService) GenerateToken() (string, error) {
	bytes := make([]byte, 64)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("unexpected signing method")
// 		}
// 		return s.secret, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	claims, ok := token.Claims.(*Claims)
// 	if !ok || !token.Valid {
// 		return nil, errors.New("invalid token")
// 	}

// 	return claims, nil
// }

// func (s *TokenService) RefreshToken(oldToken string) (string, error) {
// 	claims, err := s.ValidateToken(oldToken)
// 	if err != nil {
// 		return "", err
// 	}

// 	if time.Until(claims.ExpiresAt.Time) > time.Hour {
// 		return "", errors.New("token is still valid")
// 	}

// 	return s.GenerateToken(claims.UserID, claims.Login)
// }

// func (s *TokenService) IsTokenExpired(tokenString string) bool {
// 	claims, err := s.ValidateToken(tokenString)
// 	if err != nil {
// 		return true
// 	}

// 	return claims.ExpiresAt.Time.Before(time.Now())
// }
