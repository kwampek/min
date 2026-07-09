package handlers

import (
	"encoding/json"
	"log"
)

type WSCreatePayload struct {
	UserId int `json:"userId"`
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

func (a *API) CreatePrivateChatHandler(userID int, payload json.RawMessage) error {
	var createPayload WSCreatePayload

	if err := json.Unmarshal(payload, &createPayload); err != nil {
		log.Println("create_chat parse error:", err)
		return err
	}

	return a.createPrivateChatHandler(userID, createPayload.UserId)
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
