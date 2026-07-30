package media

import (
	"api/storage"
)

type MediaService struct {
	Storage storage.SeaweedStorage
}

func NewMediaService(Storage storage.SeaweedStorage) MediaService {
	return MediaService{
		Storage: Storage,
	}
}

func (ms *MediaService) InsertMedia(string) (int, error) {
	// TODO

	return 0, nil
}
