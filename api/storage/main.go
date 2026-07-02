package storage

import "database/sql"

type Storage struct {
	DB *sql.DB
}

func NewStorage(DB *sql.DB) *Storage {
	return &Storage{DB: DB}
}
