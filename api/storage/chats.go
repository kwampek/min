package storage

import (
	"api/models"
	"database/sql"
)

func (s *Storage) UpdateChatTitle(chatID int, title string) error {
	_, err := s.DB.Exec(`
		UPDATE Messenger.Chats
		SET title = $1
		WHERE chat_id = $2
	`, title, chatID)

	return err
}

func (s *Storage) CreateChat(typeID int, title, description string, avatarID sql.NullInt64, creatorID int) (*models.Chat, error) {
	var chat models.Chat

	err := s.DB.QueryRow(`
		INSERT INTO Messenger.Chats (
			type,
			title,
			description,
			avatar_id,
			created_at,
			creator_id
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			NOW(),
			$5
		)
		RETURNING
			chat_id,
			title,
			description,
			avatar_id,
			type,
			creator_id,
			created_at
	`, typeID, title, description, avatarID, creatorID).Scan(
		&chat.ChatID,
		&chat.Title,
		&chat.Description,
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

func (s *Storage) GetByID(chatID int) (*models.Chat, error) {
	var chat models.Chat

	err := s.DB.QueryRow(`
		SELECT
			chat_id,
			type,
			title,
			description,
			avatar_id,
			creator_id,
			created_at,
			search_privacy
		FROM Messenger.Chats
		WHERE chat_id = $1
	`, chatID).Scan(
		&chat.ChatID,
		&chat.Type,
		&chat.Title,
		&chat.Description,
		&chat.AvatarID,
		&chat.CreatorID,
		&chat.CreatedAt,
		&chat.SearchPrivacy,
	)

	if err != nil {
		return nil, err
	}

	return &chat, nil
}
