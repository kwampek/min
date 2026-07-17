package storage

func (s *Storage) GetChatMembersIDs(chatID int) ([]int, error) {
	rows, err := s.DB.Query(`
		SELECT user_id
		FROM Messenger.ChatMembers
		WHERE chat_id = $1
	`, chatID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, rows.Err()
}

func (s *Storage) AddChatMembers(chatID int, userIDs []int) ([]int, error) {
	fromIDs := make([]int, 0, len(userIDs))

	for _, uid := range userIDs {
		var fromID int
		err := s.DB.QueryRow(`
			INSERT INTO Messenger.ChatMembers (chat_id, user_id, member_from)
			VALUES ($1, $2, NOW())
			RETURNING chat_member_id
		`, chatID, uid).Scan(&fromID)
		if err != nil {
			return nil, err
		}

		fromIDs = append(fromIDs, fromID)
	}

	return fromIDs, nil
}

func (s *Storage) UpdateMemberCustomTitle(chatID, userID int, title string) error {
	_, err := s.DB.Exec(`
		UPDATE Messenger.ChatMembers
		SET custom_title = $3
		WHERE chat_id = $1
		  AND user_id = $2
	`, chatID, userID, title)

	return err
}

func (s *Storage) GetAccessRights(userID, chatID int) (int64, error) {
	var rights int64

	err := s.DB.QueryRow(`
		SELECT access_rights
		FROM Messenger.ChatMembers
		WHERE chat_id = $1
		  AND user_id = $2
	`, chatID, userID).Scan(&rights)

	return rights, err
}

func (s *Storage) HasPermission(userID, chatID int, permission int64) (bool, error) {
	var hasPermission bool
	err := s.DB.QueryRow(`
		SELECT (access_rights & $3) <> 0
		FROM Messenger.ChatMembers
		WHERE chat_id = $1
		AND user_id = $2;
	`, userID, chatID, permission).Scan(&hasPermission)

	return hasPermission, err
}
