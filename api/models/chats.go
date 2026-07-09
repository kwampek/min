package models

import (
	"database/sql"
	"time"
)

type Chat struct {
	ChatID        int            `json:"chat_id"`
	Type          int            `json:"type"`
	Title         string         `json:"title"`
	Description   sql.NullString `json:"description"`
	AvatarLink    sql.NullString `json:"avatar_link"`
	CreatorID     int            `json:"creator_id"`
	CreatedAt     time.Time      `json:"created_at"`
	SearchPrivacy sql.NullBool   `json:"search_privacy"`
}
