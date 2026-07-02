package auth

import (
	"api/models"
	"time"
)

type Storage interface {
	CreateUser(user *models.User) (int, error)
	GetUserByLogin(login string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	SaveToken(userID int, token string, expiresAt time.Time) error
	DeleteToken(userID int, token string) error
	IsTokenExists(token string) (bool, error)
}
