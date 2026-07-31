package models

import "time"

type Media struct {
	ID        int
	Media     []byte    `json:"media"`
	Preview   []byte    `json:"preview"`
	MimeType  string    `json:"mime_type"`
	SizeBytes int64     `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
}

type NewMediaRequest struct {
	Media    []byte `json:"media"`
	Preview  []byte `json:"preview"`
	MimeType string `json:"mime_type"`
}
