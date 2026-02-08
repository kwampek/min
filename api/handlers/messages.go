package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	Name         string         `json:"name"`
	AvatarLink   sql.NullString `json:"avatar"`
	ChatMemberId int            `json:"chat_member_id"`
	LastMessage  sql.NullString `json:"lastMessage"`
	Unread       int            `json:"unread"`
}

type ChatsByFoldersResponse struct {
	ChatsByFolders  map[string][]ChatsByFolder `json:"chatsByFolders"`
	MessagesByChats map[int][]Message          `json:"messagesByChats"`
}

func (A *API) loadAllMessagesHandler(w http.ResponseWriter, r *http.Request) {
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

	resp, err := A.loadAllMessages(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (A *API) loadAllMessages(userId int) (ChatsByFoldersResponse, error) {
	rowsFolders, err := A.DB.Query(`
		SELECT folder_id, folder_name
		FROM Messenger.Folders
		WHERE user_id = $1
	`, userId)
	if err != nil {
		return ChatsByFoldersResponse{}, err
	}
	defer rowsFolders.Close()

	folders := map[int]string{}
	result := ChatsByFoldersResponse{
		ChatsByFolders:  map[string][]ChatsByFolder{},
		MessagesByChats: map[int][]Message{},
	}
	result.ChatsByFolders["All"] = []ChatsByFolder{}

	for rowsFolders.Next() {
		var id int
		var name sql.NullString
		if err := rowsFolders.Scan(&id, &name); err != nil {
			return result, err
		}
		folders[id] = name.String
		result.ChatsByFolders[name.String] = []ChatsByFolder{}
	}

	rows, err := A.DB.Query(`
		WITH last_msg AS (
			SELECT DISTINCT ON (cm.chat_id)
				cm.chat_id,
				m.message_text,
				m.send_time
			FROM Messenger.ChatMembers cm
			LEFT JOIN Messenger.Messages m
				ON m.from_chat_member_id = cm.chat_member_id
			ORDER BY cm.chat_id, m.send_time DESC
		)
		SELECT 
			c.chat_id,
			c.title,
			c.avatar_link,
			cml.chat_member_id,
			COALESCE(lm.message_text, '')
		FROM Messenger.ChatMembers cml
		JOIN Messenger.Chats c ON c.chat_id = cml.chat_id
		LEFT JOIN last_msg lm ON lm.chat_id = c.chat_id
		WHERE cml.user_id = $1
	`, userId)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var chat ChatsByFolder
		if err := rows.Scan(&chat.ChatID, &chat.Name, &chat.AvatarLink, &chat.ChatMemberId, &chat.LastMessage); err != nil {
			return result, err
		}

		msgRows, err := A.DB.Query(`
			SELECT message_id, message_text, from_chat_member_id, media_id, send_time, mstatus
			FROM Messenger.Messages
			WHERE from_chat_member_id IN (
				SELECT chat_member_id FROM Messenger.ChatMembers WHERE chat_id = $1
			)
			ORDER BY send_time ASC
		`, chat.ChatID)
		if err != nil {
			return result, err
		}

		var messages = []Message{}
		for msgRows.Next() {
			var m Message
			if err := msgRows.Scan(&m.MessageID, &m.MessageText, &m.FromChatMemberID, &m.MediaID, &m.SendTime, &m.Status); err != nil {
				msgRows.Close()
				return result, err
			}
			messages = append(messages, m)
		}
		msgRows.Close()

		result.ChatsByFolders["All"] = append(result.ChatsByFolders["All"], chat)
		result.MessagesByChats[chat.ChatID] = messages

		folderRows, err := A.DB.Query(`
			SELECT folder_id FROM Messenger.FolderChats WHERE chat_id = $1
		`, chat.ChatID)

		if err != nil {
			return result, err
		}
		for folderRows.Next() {
			var fid int
			if err := folderRows.Scan(&fid); err != nil {
				continue
			}
			if folderName, ok := folders[fid]; ok {
				result.ChatsByFolders[folderName] = append(result.ChatsByFolders[folderName], chat)
			}
		}
		folderRows.Close()
	}

	return result, nil
}

func (A *API) LoadAllMessagesHandler(w http.ResponseWriter, r *http.Request) {
	A.loadAllMessagesHandler(w, r)
}

func (A *API) addMessageHandler(userId int, msg WSAddMessage) error {
	var mediaId sql.NullInt64

	if msg.Media != "" {
		log.Println("OK")
		id, err := A.insertMedia(msg.Media)
		if err != nil {
			return fmt.Errorf("media insert failed: %w", err)
		}

		mediaId = sql.NullInt64{
			Int64: int64(id),
			Valid: true,
		}
	}
	log.Println("BEFORE ADD MESSAGE")
	messageID, err := A.addMessage(msg.Text, msg.ChatMemberId, mediaId)
	if err != nil {
		return fmt.Errorf("insert message failed: %w", err)
	}

	end_msg := Message{
		MessageID:        messageID,
		MessageText:      msg.Text,
		FromChatMemberID: msg.ChatMemberId,
		MediaID:          mediaId,
		SendTime:         time.Now(),
		Status:           true,
	}

	err = A.broadcastToChat(userId, msg.ChatID, end_msg)
	return err
}

func (A *API) broadcastToChat(senderId, chatId int, message Message) error {
	allUsersFromChat, err := A.getChatMembersIds(chatId)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"type":    "new_message",
		"chat_id": chatId,
		"message": message,
	}

	for _, uid := range allUsersFromChat {
		if uid == senderId {
			continue
		}
		conn, ok := A.WsClients.Conns[uid]
		if !ok || conn == nil {
			continue
		}

		if err := conn.WriteJSON(payload); err != nil {
			log.Printf("WS send error to user %d: %v", uid, err)
			continue
		}
	}

	return nil
}

func (A *API) addMessage(text string, chatMemberId int, mediaId sql.NullInt64) (int, error) {
	log.Println("ADDMESSAGE " + text)

	var id int
	err := A.DB.QueryRow(`
        INSERT INTO Messenger.Messages
        (message_text, from_chat_member_id, media_id, send_time, mstatus)
        VALUES ($1, $2, $3, NOW(), TRUE)
        RETURNING message_id
    `, text, chatMemberId, mediaId).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
