package record

import (
	"time"
)

// UserRecord — плоская структура для маппинга строки из таблицы users.
type UserRecord struct {
	UUID         string     `db:"uuid"`
	Login        string     `db:"login"`
	PasswordHash string     `db:"password_hash"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}
