package storage

import (
	"api/models"
	"fmt"
	"strings"
)

var fields = []string{
	"login",
	"phone_number",
	"email",
	"birthday",
	"sex",
}

func (s *Storage) SearchUsers(userID int, query string) ([]models.UserWAdditionlInfo, error) {
	searchPattern := "%" + query + "%"

	rows, err := s.DB.Query(`
		SELECT 
			u.user_id, 
			u.login, 
			u.password_hash,
			u.phone_number, 
			u.email, 
			u.avatar_id, 
			u.created_at, 
			u.search_privacy,
			EXISTS (
				SELECT 1 
				FROM Messenger.Chats c
				JOIN Messenger.ChatMembers cm1 ON c.chat_id = cm1.chat_id
				JOIN Messenger.ChatMembers cm2 ON c.chat_id = cm2.chat_id
				WHERE cm1.user_id = $1
				AND cm2.user_id = u.user_id
			) AS have_chat
		FROM Messenger.Users u
		WHERE u.login ILIKE $2
		AND (u.search_privacy IS NULL OR u.search_privacy = TRUE)
		AND u.user_id != $1
		LIMIT 10
    `, userID, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserWAdditionlInfo
	for rows.Next() {
		var u models.UserWAdditionlInfo
		err := rows.Scan(
			&u.UserID,
			&u.Login,
			&u.PasswordHash,
			&u.PhoneNumber,
			&u.Email,
			&u.AvatarID,
			&u.CreatedAt,
			&u.SearchPrivacy,
			&u.HaveChat,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

func (s *Storage) UpdateProfile(userID int, payload map[string]string) error {
	args := make([]any, 0, len(fields)+1)
	set := make([]string, 0, len(fields))

	for _, field := range fields {
		if value, ok := payload[field]; ok {
			args = append(args, value)
			set = append(set, fmt.Sprintf("%s = $%d", field, len(args)))
		}
	}

	if len(set) == 0 {
		return nil
	}

	args = append(args, userID)

	query := fmt.Sprintf(`
		UPDATE Messenger.Users
		SET %s
		WHERE user_id = $%d
	`, strings.Join(set, ", "), len(args))

	_, err := s.DB.Exec(query, args...)
	return err
}
