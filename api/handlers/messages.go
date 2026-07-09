package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (A *API) LoadAllMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("user_id")
	if userIdStr == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "invalid user_id: "+userIdStr, http.StatusBadRequest)
		return
	}

	resp, err := A.MessageService.LoadAllMessages(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (A *API) AddMessageHandler(userID int, reqPayload WSAddMessagePayload) error {
	var mediaID sql.NullInt64

	if reqPayload.Media != "" {
		id, err := A.MediaService.InsertMedia(reqPayload.Media)
		if err != nil {
			return fmt.Errorf("insert media: %w", err)
		}

		mediaID = sql.NullInt64{
			Int64: int64(id),
			Valid: true,
		}
	}

	message, err := A.MessageService.AddMessage(
		reqPayload.Text,
		userID,
		reqPayload.ChatID,
		mediaID,
	)

	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"type":    "new_message",
		"chat_id": reqPayload.ChatID,
		"message": message,
	}

	users, err := A.ChatService.GetChatMembersIds(reqPayload.ChatID)
	if err != nil {
		return err
	}

	return A.WsClients.BroadcastToUsers(users, userID, payload)
}

// func (A *API) broadcastToChat(senderId, chatId int, message Message) error {
// 	allUsersFromChat, err := A.getChatMembersIds(chatId)
// 	if err != nil {
// 		return err
// 	}

// 	payload := map[string]interface{}{
// 		"type":    "new_message",
// 		"chat_id": chatId,
// 		"message": message,
// 	}

// 	for _, uid := range allUsersFromChat {
// 		if uid == senderId {
// 			continue
// 		}
// 		conn, ok := A.WsClients.Conns[uid]
// 		if !ok || conn == nil {
// 			continue
// 		}

// 		if err := conn.WriteJSON(payload); err != nil {
// 			log.Printf("WS send error to user %d: %v", uid, err)
// 			continue
// 		}
// 	}

// 	return nil
// }
