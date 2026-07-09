package handlers

import (
	"encoding/json"
	"fmt"
	"log"
)

type WSSearchPayload struct {
	Text string `json:"query"`
}

func (A *API) searchUsers(userID int, query string) error {
	users, err := A.UserService.SearchUsers(userID, query)
	if err != nil {
		return err
	}

	return A.WsService.Send(
		userID,
		map[string]any{
			"type": "search",
			"payload": map[string]any{
				"query": query,
				"users": users,
			},
		},
	)
}

func (A *API) SearchUsersHandler(userID int, payload json.RawMessage) error {
	log.Println("Nu pozya")

	var wSSearchPayload WSSearchPayload
	if err := json.Unmarshal(payload, &wSSearchPayload); err != nil {
		log.Println("search payload parse error:", err)
		return err
	}
	return A.searchUsers(userID, wSSearchPayload.Text)
}

func (A *API) EditProfile(payload map[string]string) error {
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

func (A *API) EditProfileHandler(userID int, payload json.RawMessage) error {
	var editProfilePayload map[string]string

	if err := json.Unmarshal(payload, &editProfilePayload); err != nil {
		log.Println("create folder parse error: ", err)
		return err
	}

	return A.EditProfile(editProfilePayload)
}
