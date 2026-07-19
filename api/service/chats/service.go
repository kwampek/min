package chats

import (
	"api/models"
	"database/sql"
)

type Storage interface {
	GetChatMembersIDs(chatID int) ([]int, error)
	AddChatMembers(chatID int, userIDs []int) ([]int, error)

	CreateChat(typeID int, title, description string, avatarID sql.NullInt64, creatorID int) (*models.Chat, error)
	CreatePrivateChat(userID1, userID2 int) (*models.Chat, error)

	UpdateChatTitle(chatID int, title string) error
	GetAccessRights(userID, chatID int) (int64, error)
	HasPermission(userID, chatID int, permission int64) (bool, error)
	GetByID(chatID int) (*models.Chat, error)
	UpdateMemberCustomTitle(chatID, userID int, title string) error
}

type ChatService struct {
	Storage Storage
}

func NewChatService(Storage Storage) ChatService {
	return ChatService{
		Storage: Storage,
	}
}

func (cs *ChatService) GetChatMembersIDs(chatID int) ([]int, error) {
	return cs.Storage.GetChatMembersIDs(chatID)
}

func (cs *ChatService) AddChatMembers(chatID int, userIDs []int) ([]int, error) {
	return cs.Storage.AddChatMembers(chatID, userIDs)
}

func (cs *ChatService) EditChatTitle(userID, chatID int, newName string) ([]int, error) {
	chat, err := cs.Storage.GetByID(chatID)
	if err != nil {
		return nil, err
	}

	if chat.Type == 0 {
		if err := cs.Storage.UpdateMemberCustomTitle(chatID, userID, newName); err != nil {
			return nil, err
		}

		return []int{}, nil // no need to call broadcast
	}

	hasRight, err := cs.Storage.HasPermission(userID, chatID, models.RightEditChat)
	if err != nil {
		return nil, err
	}

	if !hasRight {
		return nil, models.ErrPermissionDenied
	}

	if err := cs.Storage.UpdateChatTitle(chatID, newName); err != nil {
		return nil, err
	}

	return cs.Storage.GetChatMembersIDs(chatID)
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

func (cs *ChatService) CreateChat(typeID int, title, description string, mediaID sql.NullInt64, members []int, creatorID int) (*models.Chat, error) {
	chat, err := cs.Storage.CreateChat(typeID, title, description, mediaID, creatorID)
	if err != nil {
		return nil, err
	}

	_, err = cs.Storage.AddChatMembers(chat.ChatID, members)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (cs *ChatService) CreateGroupChannel(title, description string, mediaID sql.NullInt64, members []int, creatorID int) (*models.Chat, error) {
	return cs.CreateChat(models.ChannelChat, title, description, mediaID, members, creatorID)
}

func (cs *ChatService) CreateGroupChat(title, description string, mediaID sql.NullInt64, members []int, creatorID int) (*models.Chat, error) {
	return cs.CreateChat(models.GroupChat, title, description, mediaID, members, creatorID)
}
