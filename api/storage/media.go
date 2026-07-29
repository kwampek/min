package storage

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type SeaweedStorage struct {
	DB *sql.DB
	SC SeaweedClient
}

type SeaweedClient struct {
	MasterURL string
	VolumeURL string

	Client *http.Client
}

func New(masterURL, volumeURL string) *Client {
	return &Client{
		MasterURL: masterURL,
		VolumeURL: volumeURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *SeaweedStorage) Save(
	ctx context.Context,
	filename string,
	mimeType string,
	r io.Reader,
) (int64, error) {

	ext := filepath.Ext(filename)

	key := uuid.NewString() + ext

	size, err := s.Client.Upload(ctx, key, r)
	if err != nil {
		return 0, err
	}

	var id int64

	err = s.DB.QueryRowContext(ctx,
		`
        INSERT INTO Messenger.MediaFiles
            (filename, storage_key, mime_type, size_bytes)
        VALUES ($1,$2,$3,$4)
        RETURNING media_id
        `,
		filename,
		key,
		mimeType,
		size,
	).Scan(&id)

	return id, err
}

func (s *Storage) Open(
	ctx context.Context,
	mediaID int64,
) (io.ReadCloser, string, error) {

	var key string
	var mime string

	err := s.DB.QueryRowContext(ctx,
		`
        SELECT storage_key,mime_type
        FROM Messenger.MediaFiles
        WHERE media_id=$1
        `,
		mediaID,
	).Scan(&key, &mime)

	if err != nil {
		return nil, "", err
	}

	r, err := s.Seaweed.Open(ctx, key)

	return r, mime, err
}

func (s *Storage) Delete(
	ctx context.Context,
	mediaID int64,
) error {

	var key string

	err := s.DB.QueryRowContext(ctx,
		`
        SELECT storage_key
        FROM Messenger.MediaFiles
        WHERE media_id=$1
        `,
		mediaID,
	).Scan(&key)

	if err != nil {
		return err
	}

	if err := s.Seaweed.Delete(ctx, key); err != nil {
		return err
	}

	_, err = s.DB.ExecContext(ctx,
		`
        DELETE FROM Messenger.MediaFiles
        WHERE media_id=$1
        `,
		mediaID,
	)

	return err
}
