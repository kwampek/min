package storage

import (
	"api/models"
	"database/sql"
	"errors"
)

func (s *Storage) GetFolders(userID int) (map[int]string, error) {
	rows, err := s.DB.Query(`
		SELECT folder_id, folder_name
		FROM Messenger.Folders
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	folders := make(map[int]string)

	for rows.Next() {
		var id int
		var name sql.NullString

		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		folders[id] = name.String
	}

	return folders, rows.Err()
}

func (s *Storage) GetUserChats(userID int) ([]models.ChatsByFolder, error) {
	rows, err := s.DB.Query(`
		WITH last_msg AS (
			SELECT DISTINCT ON (cm.chat_id)
				cm.chat_id,
				m.message_text,
				m.send_time
			FROM Messenger.ChatMembers cm
			JOIN Messenger.Messages m
			ON m.from_chat_member_id = cm.chat_member_id
			ORDER BY cm.chat_id, m.send_time DESC
		)

		-- групповые
		SELECT
			c.chat_id,
			c.type,
			c.title,
			mf.file_url,
			me.chat_member_id,
			lm.message_text
		FROM Messenger.ChatMembers me
		JOIN Messenger.Chats c
			ON c.chat_id = me.chat_id
		LEFT JOIN Messenger.MediaFiles mf
			ON mf.media_id = c.avatar_id
		LEFT JOIN last_msg lm
			ON lm.chat_id = c.chat_id
		WHERE me.user_id = $1
		AND c.type <> 0

		UNION ALL

		-- приватные
		SELECT
			c.chat_id,
			c.type,
			COALESCE(other.custom_title, u.login),
			mf.file_url,
			me.chat_member_id,
			lm.message_text
		FROM Messenger.ChatMembers me
		JOIN Messenger.Chats c
			ON c.chat_id = me.chat_id
		JOIN Messenger.ChatMembers other
			ON other.chat_id = c.chat_id
		AND other.user_id <> me.user_id
		JOIN Messenger.Users u
			ON u.user_id = other.user_id
		LEFT JOIN Messenger.MediaFiles mf
			ON mf.media_id = u.avatar_id
		LEFT JOIN last_msg lm
			ON lm.chat_id = c.chat_id
		WHERE me.user_id = $1
		AND c.type = 0

		ORDER BY 6 DESC NULLS LAST;
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []models.ChatsByFolder

	for rows.Next() {
		var chat models.ChatsByFolder

		if err := rows.Scan(
			&chat.ChatID,
			&chat.Type,
			&chat.Title,
			&chat.AvatarLink,
			&chat.ChatMemberId,
			&chat.LastMessage,
		); err != nil {
			return nil, err
		}

		chats = append(chats, chat)
	}

	return chats, rows.Err()
}

func (s *Storage) GetChatMessages(chatID int) ([]models.Message, error) {
	rows, err := s.DB.Query(`
		SELECT
			message_id,
			message_text,
			from_chat_member_id,
			media_id,
			send_time,
			mstatus
		FROM Messenger.Messages
		WHERE from_chat_member_id IN (
			SELECT chat_member_id
			FROM Messenger.ChatMembers
			WHERE chat_id = $1
		)
		ORDER BY send_time ASC
	`, chatID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message

	for rows.Next() {
		var msg models.Message

		if err := rows.Scan(
			&msg.MessageID,
			&msg.MessageText,
			&msg.FromChatMemberID,
			&msg.MediaID,
			&msg.SendTime,
			&msg.Status,
		); err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (s *Storage) GetChatFolderIDs(chatID int) ([]int, error) {
	rows, err := s.DB.Query(`
		SELECT folder_id
		FROM Messenger.FolderChats
		WHERE chat_id = $1
	`, chatID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int

	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (s *Storage) AddMessage(text string, userID int, chatID int, mediaID sql.NullInt64) (messageID int, chatMemberID int, err error) {
	err = s.DB.QueryRow(`
		WITH member AS (
			SELECT chat_member_id
			FROM Messenger.ChatMembers
			WHERE user_id = $2
			  AND chat_id = $3
		),
		inserted AS (
			INSERT INTO Messenger.Messages
			(
				message_text,
				from_chat_member_id,
				media_id,
				send_time,
				mstatus
			)
			SELECT
				$1,
				chat_member_id,
				$4,
				NOW(),
				TRUE
			FROM member
			RETURNING message_id, from_chat_member_id
		)
		SELECT message_id, from_chat_member_id
		FROM inserted;
	`,
		text,
		userID,
		chatID,
		mediaID,
	).Scan(&messageID, &chatMemberID)

	if err == sql.ErrNoRows {
		return 0, 0, errors.New("user is not member of chat")
	}

	return
}
