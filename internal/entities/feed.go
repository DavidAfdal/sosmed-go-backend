package entities

import (
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Caption   string    `db:"caption"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type FeedMedia struct {
	ID   uuid.UUID `db:"id"`
	Url  string    `db:"id"`
	Type string    `db:"id"`
}
