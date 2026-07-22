package models

import (
	"database/sql"
	"errors"
	"time"
)

const (
	RightEditChat       = 1 << iota // 1
	RightDeleteMessages             // 2
	RightBanUsers                   // 4
	RightInviteUsers                // 8
	RightPinMessages                // 16
	RightEditInfo                   // 32
)

const (
	PrivateChat = 0
	GroupChat   = 1
	ChannelChat = 2
)

var (
	ErrPermissionDenied = errors.New("permission denied")
)

type Chat struct {
	ChatID        int            `json:"chat_id" db:"chat_id"`
	Type          int            `json:"type" db:"type"`
	Title         string         `json:"title" db:"title"`
	Description   sql.NullString `json:"description" db:"description"`
	AvatarID      sql.NullInt64  `json:"avatar_id" db:"avatar_id"`
	CreatorID     int            `json:"creator_id" db:"creator_id"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	SearchPrivacy sql.NullBool   `json:"search_privacy" db:"search_privacy"`
}
