package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Avatar    string    `db:"avatar"`
	Bio       string    `db:"bio"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
