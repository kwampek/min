package models

import (
	"database/sql"
	"time"
)

type Message struct {
	MessageID        int           `json:"message_id"`
	MessageText      string        `json:"message_text"`
	FromChatMemberID int           `json:"from_chat_member_id"`
	MediaID          sql.NullInt64 `json:"media_id"`
	SendTime         time.Time     `json:"send_time"`
	Status           bool          `json:"status"`
}

type ChatsByFolder struct {
	ChatID       int            `json:"chat_id"`
	Title        string         `json:"title"`
	AvatarLink   sql.NullString `json:"avatar"`
	ChatMemberId int            `json:"chat_member_id"`
	LastMessage  sql.NullString `json:"lastMessage"`
	Unread       int            `json:"unread"`
}

type ChatsByFoldersResponse struct {
	ChatsByFolders  map[string][]ChatsByFolder `json:"chatsByFolders"`
	MessagesByChats map[int][]Message          `json:"messagesByChats"`
}
