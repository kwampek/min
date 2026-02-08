package handlers

import (
	"fmt"
)

func (A *API) createFolderHandler(user_id int, name string) (int, error) {
	var folder_id int
	err := A.DB.QueryRow(`
        INSERT INTO Messenger.Folders (user_id, folder_name)
        VALUES ($1, $2)
        RETURNING folder_id
    `, user_id, name).Scan(&folder_id)
	if err != nil {
		return 0, err
	}
	return folder_id, nil
}

func (A *API) editFolderNameHandler(user_id int, prev_name, name string) error {
	res, err := A.DB.Exec(`
        UPDATE Messenger.Folders
        SET folder_name = $1
        WHERE user_id = $2 and folder_name = $3
    `, name, user_id, prev_name)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("folder with name %s not found", prev_name)
	}

	return nil
}

func (A *API) addChatToFolderHandler(user_id, chat_id int, folder_name string) error {
	_, err := A.DB.Exec(`
		INSERT INTO Messenger.FolderChats (folder_id, chat_id)
		SELECT folder_id, $1
		FROM Messenger.Folders
		WHERE user_id = $2 AND folder_name = $3
	`, chat_id, user_id, folder_name)

	return err
}

func (A *API) removeChatFromFolderHandler(user_id, chat_id int, folder_name string) error {
	_, err := A.DB.Exec(`
		DELETE FROM Messenger.FolderChats
		WHERE chat_id = $1
		AND folder_id IN (
			SELECT folder_id
			FROM Messenger.Folders
			WHERE user_id = $2 AND folder_name = $3
		)
	`, chat_id, user_id, folder_name)

	return err
}
