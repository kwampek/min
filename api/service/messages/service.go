package messages

import (
	"api/models"
	"database/sql"
	"errors"
	"time"
)

type Storage interface {
	AddMessage(text string, userID, chatID int, mediaId sql.NullInt64) (int, int, error)
	GetFolders(userID int) (map[int]string, error)
	GetUserChats(userID int) ([]models.ChatsByFolder, error)
	GetChatMessages(chatID int) ([]models.Message, error)
	GetChatFolderIDs(chatID int) ([]int, error)
}

type MessageService struct {
	Storage Storage
}

func NewMessagesService(storage Storage) MessageService {
	return MessageService{
		Storage: storage,
	}
}

var (
	ErrAddMessage       = errors.New("internal error")
	ErrBroadcastFailure = errors.New("internal error")
)

func (ms *MessageService) LoadAllMessages(userID int) (models.ChatsByFoldersResponse, error) {

	folders, err := ms.Storage.GetFolders(userID)
	if err != nil {
		return models.ChatsByFoldersResponse{}, err
	}

	chats, err := ms.Storage.GetUserChats(userID)
	if err != nil {
		return models.ChatsByFoldersResponse{}, err
	}

	result := models.ChatsByFoldersResponse{
		ChatsByFolders:  make(map[string][]models.ChatsByFolder),
		MessagesByChats: make(map[int][]models.Message),
	}

	result.ChatsByFolders["All"] = []models.ChatsByFolder{}

	for _, folderName := range folders {
		result.ChatsByFolders[folderName] = []models.ChatsByFolder{}
	}

	for _, chat := range chats {

		messages, err := ms.Storage.GetChatMessages(chat.ChatID)
		if err != nil {
			return result, err
		}

		result.MessagesByChats[chat.ChatID] = messages
		result.ChatsByFolders["All"] = append(result.ChatsByFolders["All"], chat)

		folderIDs, err := ms.Storage.GetChatFolderIDs(chat.ChatID)
		if err != nil {
			return result, err
		}

		for _, id := range folderIDs {
			if folderName, ok := folders[id]; ok {
				result.ChatsByFolders[folderName] = append(result.ChatsByFolders[folderName], chat)
			}
		}
	}

	return result, nil
}

func (ms *MessageService) AddMessage(text string, userID int, chatID int, mediaID sql.NullInt64) (models.Message, error) {
	messageID, chatMemberID, err := ms.Storage.AddMessage(text, userID, chatID, mediaID)
	if err != nil {
		return models.Message{}, err
	}

	msg := models.Message{
		MessageID:        messageID,
		MessageText:      text,
		FromChatMemberID: chatMemberID,
		MediaID:          mediaID,
		SendTime:         time.Now(),
		Status:           true,
	}

	return msg, nil
}
