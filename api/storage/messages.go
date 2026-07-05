package storage

import (
	"database/sql"
	"log"
)

func (s *Storage) AddMessage(text string, chatMemberId int, mediaId sql.NullInt64) (int, error) {
	log.Println("ADDMESSAGE " + text)

	var id int
	err := s.DB.QueryRow(`
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
