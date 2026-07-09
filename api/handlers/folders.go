package handlers

import (
	"encoding/json"
	"log"
)

type WSCreateFolderPayload struct {
	Name string `json:"folder_name"`
}

func (A *API) CreateFolderHandler(userID int, payload json.RawMessage) error {
	var createFolderPayload WSCreateFolderPayload

	if err := json.Unmarshal(payload, &createFolderPayload); err != nil {
		log.Println("create folder parse error: ", err)
		return err
	}

	_, err := A.FolderService.CreateFolder(userID, createFolderPayload.Name)
	// TODO - check necessity of folderID return
	return err
}

type WSEditFolderNamePayload struct {
	PrevName string `json:"prev_name"`
	NewName  string `json:"new_name"`
}

func (A *API) EditFolderNameHandler(userID int, payload json.RawMessage) error {
	var editFolderNamePayload WSEditFolderNamePayload

	if err := json.Unmarshal(payload, &editFolderNamePayload); err != nil {
		log.Println("create folder parse error: ", err)
		return err
	}

	return A.FolderService.EditFolderName(userID, editFolderNamePayload.PrevName, editFolderNamePayload.NewName)

}

type WSAddChatToFolderPayload struct {
	ChatId     int    `json:"chat_id"`
	Foldername string `json:"folder_name"`
}

func (A *API) AddChatToFolderHandler(userID int, payload json.RawMessage) error {
	var AddChatToFolderPayload WSAddChatToFolderPayload

	if err := json.Unmarshal(payload, &AddChatToFolderPayload); err != nil {
		log.Println("create folder parse error: ", err)
		return err
	}

	return A.FolderService.AddChatToFolder(userID, AddChatToFolderPayload.ChatId, AddChatToFolderPayload.Foldername)
}

type WSRemoveChatFromFolderPayload struct {
	ChatId     int    `json:"chat_id"`
	Foldername string `json:"folder_name"`
}

func (A *API) RemoveChatFromFolderHandler(userID int, payload json.RawMessage) error {
	var RemoveChatFromFolderPayload WSRemoveChatFromFolderPayload

	if err := json.Unmarshal(payload, &RemoveChatFromFolderPayload); err != nil {
		log.Println("remove chat from folder parse error: ", err)
		return err
	}

	return A.FolderService.RemoveChatFromFolder(userID, RemoveChatFromFolderPayload.ChatId, RemoveChatFromFolderPayload.Foldername)
}

type WSToggleChatInFolderPayload struct {
	ChatId     int    `json:"chat_id"`
	Foldername string `json:"folder_name"`
	IsChecked  bool   `json:"is_checked"`
}

func (A *API) ToggleChatInFolderHandler(userID int, payload json.RawMessage) error {
	var ToggleChatInFolderPayload WSToggleChatInFolderPayload

	if err := json.Unmarshal(payload, &ToggleChatInFolderPayload); err != nil {
		log.Println("toggle chat in folder parse error: ", err)
		return err
	}

	if ToggleChatInFolderPayload.IsChecked {
		return A.FolderService.AddChatToFolder(userID, ToggleChatInFolderPayload.ChatId, ToggleChatInFolderPayload.Foldername)
	}

	return A.FolderService.RemoveChatFromFolder(userID, ToggleChatInFolderPayload.ChatId, ToggleChatInFolderPayload.Foldername)

}
