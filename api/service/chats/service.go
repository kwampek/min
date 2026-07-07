package chats

type Storage interface {
	GetChatMembersIds(chatID int) ([]int, error)
	AddChatMembers(chatID int, userIDs []int) ([]int, error)
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
