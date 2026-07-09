package chats

import "api/models"

type Storage interface {
	GetChatMembersIds(chatID int) ([]int, error)
	AddChatMembers(chatID int, userIDs []int) ([]int, error)

	CreatePrivateChat(userID1, userID2 int) (*models.Chat, error)
	UpdateChatTitle(chatID int, title string) error
}

type ChatService struct {
	Storage Storage
}

func NewChatService(Storage Storage) ChatService {
	return ChatService{
		Storage: Storage,
	}
}

func (cs *ChatService) GetChatMembersIds(chatID int) ([]int, error) {
	return cs.Storage.GetChatMembersIds(chatID)
}

func (cs *ChatService) AddChatMembers(chatID int, userIDs []int) ([]int, error) {
	return cs.Storage.AddChatMembers(chatID, userIDs)
}

func (cs *ChatService) EditChatTitle(chatID int, title string) error {
	return cs.Storage.UpdateChatTitle(chatID, title)
}

func (cs *ChatService) CreatePrivateChat(userID1, userID2 int) (*models.Chat, []int, error) {
	chat, err := cs.Storage.CreatePrivateChat(userID1, userID2)
	if err != nil {
		return nil, nil, err
	}

	members := []int{
		userID1,
		userID2,
	}

	_, err = cs.Storage.AddChatMembers(chat.ChatID, members)
	if err != nil {
		return nil, nil, err
	}

	return chat, members, nil
}
