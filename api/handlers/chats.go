package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
)

type WSCreatePayload struct {
	UserId int `json:"userId"`
}

type WSCreateChatGroupPayload struct {
	Type        int
	Title       string
	Description string
	AvatarID    sql.NullInt64
	Users       []int `json:"users"`
}

func (a *API) createPrivateChatHandler(userID1, userID2 int) error {
	chat, members, err := a.ChatService.CreatePrivateChat(userID1, userID2)
	if err != nil {
		return err
	}

	return a.WsService.BroadcastToUsers(
		members,
		-1,
		map[string]any{
			"type":    "create_chat",
			"payload": chat,
		},
	)
}

func (a *API) createChatHandler(typeID int, title, description string, avatarID sql.NullInt64, creatorID int, members []int) error {
	chat, err := a.ChatService.CreateChat(typeID, title, description, avatarID, members, creatorID)
	if err != nil {
		return err
	}

	// TODO clear -1 here

	return a.WsService.BroadcastToUsers(
		members,
		-1,
		map[string]any{
			"type":    "create_chat",
			"payload": chat,
		},
	)
}

func (a *API) CreatePrivateChatHandler(userID int, payload json.RawMessage) error {
	var createPayload WSCreatePayload

	if err := json.Unmarshal(payload, &createPayload); err != nil {
		log.Println("create_chat parse error:", err)
		return err
	}

	return a.createPrivateChatHandler(userID, createPayload.UserId)
}

func (a *API) CreateChatHandler(userID int, payload json.RawMessage) error {
	var createCGPayload WSCreateChatGroupPayload

	if err := json.Unmarshal(payload, &createCGPayload); err != nil {
		log.Println("create_chat parse error:", err)
		return err
	}

	return a.createChatHandler(
		createCGPayload.Type,
		createCGPayload.Title,
		createCGPayload.Description,
		createCGPayload.AvatarID,
		userID,
		createCGPayload.Users,
	)
}

type WSEditChatNamePayload struct {
	ChatId  int    `json:"chat_id"`
	NewName string `json:"new_name"`
}

func (a *API) EditChatNameHandler(userID int, payload json.RawMessage) error {
	var editChatNamePayload WSEditChatNamePayload

	if err := json.Unmarshal(payload, &editChatNamePayload); err != nil {
		log.Println("create folder parse error: ", err)
		return err
	}

	return a.ChatService.EditChatTitle(editChatNamePayload.ChatId, editChatNamePayload.NewName)
}
