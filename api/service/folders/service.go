package folders

type Storage interface {
	CreateFolder(userID int, name string) (int, error)
	EditFolderName(userID int, prev_name, name string) error

	AddChatToFolder(userID, chatID int, folderName string) error
	RemoveChatFromFolder(userID, chatID int, folderName string) error
}

type FolderService struct {
	Storage Storage
}

func NewFolderService(Storage Storage) FolderService {
	return FolderService{
		Storage: Storage,
	}
}

func (fs *FolderService) CreateFolder(userID int, name string) (int, error) {
	return fs.Storage.CreateFolder(userID, name)
}

func (fs *FolderService) EditFolderName(userID int, prev_name, name string) error {
	return fs.Storage.EditFolderName(userID, prev_name, name)
}

func (fs *FolderService) AddChatToFolder(userID, chatID int, folderName string) error {
	return fs.Storage.AddChatToFolder(userID, chatID, folderName)
}

func (fs *FolderService) RemoveChatFromFolder(userID, chatID int, folderName string) error {
	return fs.Storage.RemoveChatFromFolder(userID, chatID, folderName)
}
