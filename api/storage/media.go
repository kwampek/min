package storage

import (
	"api/models"
	"database/sql"
	"fmt"
)

func (s *Storage) CreateMedia(media, preview []byte, mime string) (int, error) {
	var id int

	err := s.DB.QueryRow(`
        INSERT INTO Messenger.MediaFiles
            (media, preview, mime_type, size_bytes)
        VALUES
            ($1, $2, $3, $4)
        RETURNING media_id
    `,
		media,
		preview,
		mime,
		len(media),
	).Scan(&id)

	return id, err
}

func (s *Storage) GetImage(id int, getPreview bool) ([]byte, string, error) {
	var (
		data []byte
		mime string
	)

	var fieldName string
	if getPreview {
		fieldName = "preview"
	} else {
		fieldName = "media"
	}

	query := fmt.Sprintf(`
        SELECT %s, mime_type
        FROM Messenger.MediaFiles
        WHERE media_id = $1
    `, fieldName)

	err := s.DB.QueryRow(query, id).Scan(&data, &mime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", fmt.Errorf("media with id %d not found", id)
		}
		return nil, "", fmt.Errorf("failed to get image: %w", err)
	}

	if data == nil {
		if getPreview {
			return nil, "", fmt.Errorf("preview not available for media %d", id)
		}
		return nil, "", fmt.Errorf("full media not available for media %d", id)
	}

	return data, mime, nil
}

func (s *Storage) GetMediaWithPreview(id int) (*models.Media, error) {
	var m models.Media

	err := s.DB.QueryRow(`
        SELECT
            media_id,
            media,
            preview,
            mime_type,
            size_bytes,
            created_at
        FROM Messenger.MediaFiles
        WHERE media_id = $1
    `, id).Scan(
		&m.ID,
		&m.Media,
		&m.Preview,
		&m.MimeType,
		&m.SizeBytes,
		&m.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *Storage) DeleteMedia(id int) error {
	result, err := s.DB.Exec(`
        DELETE FROM Messenger.MediaFiles
        WHERE media_id = $1
    `, id)

	if err != nil {
		return fmt.Errorf("failed to hard delete media: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("media with id %d not found", id)
	}

	return nil
}
