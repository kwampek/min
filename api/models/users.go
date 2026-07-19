package models

import (
	"database/sql"
	"time"
)

type User struct {
	UserID        int           `json:"user_id"`
	Login         string        `json:"name"`
	PasswordHash  string        `json:"password"`
	PhoneNumber   string        `json:"phone_number"`
	Email         string        `json:"email"`
	AvatarID      sql.NullInt32 `json:"avatar_id"`
	CreatedAt     time.Time     `json:"created_at"`
	SearchPrivacy bool          `json:"search_privacy"`
}

// TODO remove it

type UserWAdditionlInfo struct {
	UserID        int           `json:"user_id"`
	Login         string        `json:"name"`
	PasswordHash  string        `json:"password"`
	PhoneNumber   string        `json:"phone_number"`
	Email         string        `json:"email"`
	AvatarID      sql.NullInt32 `json:"avatar_id"`
	CreatedAt     time.Time     `json:"created_at"`
	SearchPrivacy bool          `json:"search_privacy"`
	HaveChat      bool          `json:"have_chat"` // only for SearchUsersHandler
}
