package storage

import (
	"api/models"
)

func (s *Storage) UpdateChatTitle(chatID int, title string) error {
	_, err := s.DB.Exec(`
		UPDATE Messenger.Chats
		SET title = $1
		WHERE chat_id = $2
	`, title, chatID)

	return err
}

func (s *Storage) CreateChat(title string, creatorID int, typeID int) (*models.Chat, error) {
	var chat models.Chat

	err := s.DB.QueryRow(`
		INSERT INTO Messenger.Chats (
			type,
			title,
			created_at,
			creator_id
		)
		VALUES (
			$1,
			$2
			NOW(),
			$3
		)
		RETURNING
			chat_id,
			title,
			type,
			creator_id,
			created_at
	`, typeID, title, creatorID).Scan(
		&chat.ChatID,
		&chat.Title,
		&chat.Type,
		&chat.CreatorID,
		&chat.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (s *Storage) CreatePrivateChat(userID1, userID2 int) (*models.Chat, error) {
	var chat models.Chat

	err := s.DB.QueryRow(`
		INSERT INTO Messenger.Chats (
			type,
			title,
			created_at,
			creator_id
		)
		VALUES (
			1,
			(
				SELECT u1.login || '-' || u2.login
				FROM Messenger.Users u1,
					 Messenger.Users u2
				WHERE u1.user_id = $1
				  AND u2.user_id = $2
			),
			NOW(),
			$1
		)
		RETURNING
			chat_id,
			title,
			type,
			creator_id,
			created_at
	`, userID1, userID2).Scan(
		&chat.ChatID,
		&chat.Title,
		&chat.Type,
		&chat.CreatorID,
		&chat.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &chat, nil
}
