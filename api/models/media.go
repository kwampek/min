package models

import (
	"database/sql"
)

type Media struct {
	MediaID sql.NullInt64
	Link    string
}
