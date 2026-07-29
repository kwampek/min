package storage

import (
	"database/sql"
	"io"

	"api/util/seaweed"
)

type SeaweedStorage struct {
	DB *sql.DB
	SC seaweed.SeaweedClient
}

func MakeSeaweedStorage(DB *sql.DB, masterURL, volumeURL string) SeaweedStorage {
	return SeaweedStorage{
		DB: DB,
		SC: seaweed.MakeSeaweedClient(masterURL, volumeURL),
	}
}

func (s *SeaweedStorage) Save(filename string, mimeType string, r io.Reader) (int64, error) {

	key, size, err := s.SC.Upload(r, filename)
	if err != nil {
		return 0, err
	}

	var id int64

	err = s.DB.QueryRow(
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

func (s *SeaweedStorage) Open(mediaID int64) (io.ReadCloser, string, error) {

	var key string
	var mime string

	err := s.DB.QueryRow(
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

	r, err := s.SC.Open(key)

	return r, mime, err
}

func (s *SeaweedStorage) Delete(mediaID int64) error {
	var key string

	err := s.DB.QueryRow(
		`
        DELETE FROM Messenger.MediaFiles
        WHERE media_id=$1
		RETURNING storage_key
        `,
		mediaID,
	).Scan(&key)

	if err != nil {
		return err
	}

	if err := s.SC.Delete(key); err != nil {
		return err
	}

	return err
}
