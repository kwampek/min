package storage

func (s *Storage) getChatMembersIds(chatID int) ([]int, error) {
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

func (s *Storage) addChatMembers(chatID int, userIDs []int) ([]int, error) {
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
