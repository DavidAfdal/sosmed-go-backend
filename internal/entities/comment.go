package entities

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	FeedID    uuid.UUID `db:"feed_id"`
	ParentID  uuid.UUID `db:"parent_id"`
	Comment   string    `db:"comment"`
	ReplyCout int       `db:"reply_count"`
	CreatedAt time.Time `db:"created_at"`
	User      *User
}
