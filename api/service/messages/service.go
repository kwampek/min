package messages

type MessageService struct {
	Storage Storage
}

func NewService(storage Storage) MessageService {
	return MessageService{
		Storage: storage,
	}
}
