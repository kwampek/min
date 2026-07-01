package models

import (
	"time"
)

type User struct {
	UserID        int       `json:"user_id"`
	Login         string    `json:"name"`
	PasswordHash  string    `json:"password"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	AvatarLink    string    `json:"avatar"`
	CreatedAt     time.Time `json:"created_at"`
	SearchPrivacy bool      `json:"search_privacy"`
	// TODO	HaveChat      bool           `json:"have_chat"` // only for SearchUsersHandler
}
