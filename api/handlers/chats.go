package handlers

import (
	"database/sql"
	"log"
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

type SChat struct {
	ChatID        int            `json:"chat_id"`
	Type          int            `json:"type"`
	Title         string         `json:"name"`
	Description   sql.NullString `json:"description"`
	AvatarLink    sql.NullString `json:"avatar_link"`
	CreatorID     int            `json:"creator_id"`
	CreatedAt     time.Time      `json:"created_at"`
	SearchPrivacy sql.NullBool   `json:"search_privacy"`
	ChatMemberID  int            `json:"chat_member_id"`
}

func (a *API) editChatNameHandler(chat_id int, title string) error {
	_, err := a.DB.Exec(`
        UPDATE Messenger.Chats
        SET title=$1
        WHERE chat_id=$2`,
		title, chat_id)

	return err
}

func (A *API) createPrivateChatHandler(userID1, userID2 int) error {
	var chatID int
	var chatTitle string

	err := A.DB.QueryRow(`
        INSERT INTO Messenger.Chats (type, title, created_at, creator_id)
    	VALUES (
			1,
			(SELECT u1.login || '-' || u2.login
			FROM Messenger.Users u1, Messenger.Users u2
			WHERE u1.user_id = $1 AND u2.user_id = $2),
			NOW(),
			$1
		)
		RETURNING chat_id, title
    `, userID1, userID2).Scan(&chatID, &chatTitle) // may be can send login in query

	if err != nil {
		log.Println("ERORR: ", err)
		return err
	}

	busers := []int{userID1, userID2}

	fcm, err := A.addChatMembers(chatID, busers)
	if err != nil {
		log.Println("AddChatMembers err", err)
		return err
	}

	var chat = SChat{ChatID: chatID, Title: chatTitle, CreatorID: userID1, Type: 1}

	for i, uid := range busers {
		conn, ok := A.WsClients.Conns[uid]
		if !ok || conn == nil {
			log.Println("Conn err", err)
			continue
		}

		chat.ChatMemberID = fcm[i]

		wsPayload := map[string]interface{}{
			"type":    "create_chat",
			"payload": chat,
		}
		if err := conn.WriteJSON(wsPayload); err != nil {
			log.Printf("WS send error to user %d: %v", uid, err)
		}
	}

	return nil
}
