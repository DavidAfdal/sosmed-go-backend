package entities

import (
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Caption   string    `db:"caption"`
	User      *User
	Medias    []*FeedMedia
	Likes     int
	Comments  int
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type FeedMedia struct {
	ID     uuid.UUID `db:"id"`
	FeedId uuid.UUID `db:"feed_id"`
	Url    string    `db:"id"`
	Type   string    `db:"id"`
}
