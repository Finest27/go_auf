package models

import (
	"database/sql"
)

type Feed struct {
	ID        int64
	URL       string
	Name      string
	LastCheck sql.NullTime
	Active    bool
}
