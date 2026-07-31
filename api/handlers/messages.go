package handlers

import (
	"api/models"
	"database/sql"
	"encoding/json"
	"log"
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

type WSAddMessagePayload struct {
	ChatID       int                    `json:"chat_id"`
	ChatMemberId int                    `json:"from_chat_member_id"`
	Text         string                 `json:"text"`
	MediaID      sql.NullInt64          `json:"media_id"`
	Media        models.NewMediaRequest `json:"media"`
}

func (A *API) GetMediaIDFromMessage(mediaID sql.NullInt64, media models.NewMediaRequest) (sql.NullInt64, error) {
	if !mediaID.Valid {
		return mediaID, nil
	}

	if val, _ := mediaID.Value(); val.(int) != -1 {
		return mediaID, nil
	}

	id, err := A.MediaService.CreateMedia(
		media.Media,
		media.Preview,
		media.MimeType,
	)

	return sql.NullInt64{
		Int64: int64(id),
		Valid: true,
	}, err

}

func (A *API) addMessageHandler(userID int, reqPayload WSAddMessagePayload) error {
	mediaID, err := A.GetMediaIDFromMessage(reqPayload.MediaID, reqPayload.Media)
	if err != nil {
		return err
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

	users, err := A.ChatService.GetChatMembersIDs(reqPayload.ChatID)
	if err != nil {
		return err
	}

	return A.WsService.BroadcastToUsers(users, userID, payload)
}

func (A *API) AddMessageHandler(ID models.Identifier, payload json.RawMessage) error {
	var msg WSAddMessagePayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		log.Println("new_message parse error:", err)
		return err
	}
	return A.addMessageHandler(ID.UserID, msg)
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
