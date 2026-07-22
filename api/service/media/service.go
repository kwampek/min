package media

import "log"

type Storage interface {
}

type MediaService struct {
	Storage Storage
}

func NewMediaService(Storage Storage) MediaService {
	return MediaService{
		Storage: Storage,
	}
}
func (ms *MediaService) InsertMedia(string) (int, error) {
	// TODO
	log.Panic("UNIMPLEMENTED")
	return 0, nil
}
