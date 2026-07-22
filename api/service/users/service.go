package users

import (
	"api/models"
)

type Storage interface {
	SearchUsers(userID int, query string) ([]models.UserWAdditionlInfo, error)
	UpdateProfile(userID int, payload map[string]string) error
}
type UserService struct {
	Storage Storage
}

func NewUserService(Storage Storage) UserService {
	return UserService{
		Storage: Storage,
	}
}

func (us *UserService) SearchUsers(userID int, query string) ([]models.UserWAdditionlInfo, error) {
	return us.Storage.SearchUsers(userID, query)
}
