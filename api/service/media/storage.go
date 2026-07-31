package media

import "api/models"

type Storage interface {
	CreateMedia(media, preview []byte, mimeType string) (int, error)

	GetImage(id int, getPreview bool) ([]byte, string, error)
	GetMediaWithPreview(id int) (*models.Media, error)

	DeleteMedia(id int) error
}
