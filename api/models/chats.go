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
	PrivateChatType = 0
	GroupChatType
	GroupChanellType
)

var (
	ErrPermissionDenied = errors.New("permission denied")
)

type Chat struct {
	ChatID        int            `json:"chat_id"`
	Type          int            `json:"type"`
	Title         string         `json:"title"`
	Description   sql.NullString `json:"description"`
	AvatarID      sql.NullInt64  `json:"avatar_link"`
	CreatorID     int            `json:"creator_id"`
	CreatedAt     time.Time      `json:"created_at"`
	SearchPrivacy sql.NullBool   `json:"search_privacy"`
}
