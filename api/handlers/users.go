package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type User struct {
	UserID        int            `json:"user_id"`
	Login         string         `json:"name"` //login
	Password      string         `json:"password"`
	PhoneNumber   sql.NullString `json:"phone_number"`
	AvatarLink    sql.NullString `json:"avatar"` //avatar_link
	CreatedAt     time.Time      `json:"created_at"`
	SearchPrivacy sql.NullBool   `json:"search_privacy"`
	HaveChat      bool           `json:"have_chat"` // only for SearchUsersHandler
}

/*
type WSEditProfilePayload struct {
	UserId      int    `json:"user_id"`
	Login       string `json:"login"`
	Email string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Birthday    string `json:"birthday"`
	Sex         string `json:"sex"`
}
*/

func (A *API) editProfileHandler(payload map[string]string) error {
	user_id, ok := payload["user_id"]
	if !ok {
		return fmt.Errorf("missing user_id")
	}

	fields := []string{"login, phone_number, email, birthday, sex"}
	args := []string{}

	query := "UPDATE Messenger.Users SET "

	for _, arg := range fields {
		val, ok := payload[arg]
		if ok {
			args = append(args, val)
			query += fmt.Sprintf("%s = $%d,", arg, len(args))
		}
	}

	if len(args) == 0 {
		return nil
	}

	query += fmt.Sprintf(" WHERE user_id=%d", user_id)

	_, err := A.DB.Exec(query, args)
	return err
}

func (A *API) searchUsersHandler(userId int, query string) error {
	users, err := A.searchUsers(userId, query)

	if err != nil {
		log.Println("Error while searching users: ", err)
		return err
	}

	payload := map[string]interface{}{
		"type": "search",
		"payload": map[string]interface{}{
			"query": query,
			"users": users,
		},
	}

	conn, ok := A.WsClients.Conns[userId]
	if !ok || conn == nil {
		log.Println("No ws connection for user", userId)
		return fmt.Errorf("no ws connection for user %d", userId)
	}

	if err := conn.WriteJSON(payload); err != nil {
		log.Println("Ws Send error to user ", userId, err)
		return fmt.Errorf("WS send error to user %d: %v", userId, err)
	}

	return nil
}

func (A *API) searchUsers(userId int, query string) ([]User, error) {
	searchPattern := "%" + query + "%"

	rows, err := A.DB.Query(`
		SELECT 
			u.user_id, 
			u.login, 
			u.password, 
			u.phone_number, 
			u.avatar_link, 
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
    `, userId, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.UserID,
			&u.Login,
			&u.Password,
			&u.PhoneNumber,
			&u.AvatarLink,
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

func (A *API) registerUser(login, password string) (int, error) {
	var user_id int

	err := A.DB.QueryRow(`
		INSERT INTO Messenger.Users (login, password, created_at)
		VALUES ($1, $2, NOW())
		RETURNING user_id`, login, password).Scan(&user_id)

	if err != nil {
		return 0, err
	}

	return user_id, nil
}
