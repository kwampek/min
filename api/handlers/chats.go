package handlers

import (
	"api/models"
	"database/sql"
	"encoding/json"
	"log"
)

type WSCreatePayload struct {
	Type        int             `json:"type"`
	Title       string          `json:"title"`
	Description string          `json:"desc"`
	Avatar      WSAvatarPayload `json:"avatar"`
	Users       []int           `json:"users"`
}

type WSAvatarPayload struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	Data string `json:"data"`
}

func (a *API) createChatHandler(typeID int, title, description string, avatarID sql.NullInt64, creatorID int, members []int) error {
	chat, err := a.ChatService.CreateChat(typeID, title, description, avatarID, members, creatorID)
	if err != nil {
		return err
	}

	// TODO clear -1 here
	// sobad
	members = append(members, creatorID)

	return a.WsService.BroadcastToUsers(
		members,
		-1,
		map[string]any{
			"type":    "create_chat",
			"payload": chat,
		},
	)
}

func (a *API) CreateChatHandler(ID models.Identifier, payload json.RawMessage) error {
	var createCGPayload WSCreatePayload

	if err := json.Unmarshal(payload, &createCGPayload); err != nil {
		log.Println("create_chat parse error:", err)
		return err
	}

	// Add Avatar

	return a.createChatHandler(
		createCGPayload.Type,
		createCGPayload.Title,
		createCGPayload.Description,
		sql.NullInt64{},
		ID.UserID,
		createCGPayload.Users,
	)
}

type WSEditChatNamePayload struct {
	ChatId  int    `json:"chat_id"`
	NewName string `json:"new_name"`
}

type WSChatNameChanged struct {
	ChatID  int    `json:"chat_id"`
	NewName string `json:"new_name"`
}

func (a *API) EditChatNameHandler(ID models.Identifier, payload json.RawMessage) error {
	var p WSEditChatNamePayload

	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}

	memberIDs, err := a.ChatService.EditChatTitle(ID.UserID, p.ChatId, p.NewName)
	if err != nil {
		return err
	}

	return a.WsService.BroadcastToUsers(
		memberIDs,
		-1,
		WSChatNameChanged{
			ChatID:  p.ChatId,
			NewName: p.NewName,
		},
	)
}
