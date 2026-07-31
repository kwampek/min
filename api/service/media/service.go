package media

import "api/models"

type MediaService struct {
	Storage Storage
}

func NewMediaService(storage Storage) MediaService {
	return MediaService{
		Storage: storage,
	}
}

func (s *MediaService) CreateMedia(media, preview []byte, mime string) (int, error) {
	return s.Storage.CreateMedia(media, preview, mime)
}

func (s *MediaService) GetMediaWithPreview(id int) (*models.Media, error) {
	return s.Storage.GetMediaWithPreview(id)
}

func (s *MediaService) DeleteMedia(id int) error {
	return s.Storage.DeleteMedia(id)
}
