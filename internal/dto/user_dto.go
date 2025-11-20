package dto

import "time"

type CreateUserRequest struct {
	Username string `form:"username" validate:"required"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
	Avatar   string `form:"avatar"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdatedUserRequest struct {
	Username string `form:"username" validate:"required"`
	Avatar   string `form:"avatar"`
	Bio      string `form:"bio"`
}

type UserResponse struct {
	ID        string    `json:"id,omitzero"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitzero"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio,omitzero"`
	Followers int       `json:"followers,omitzero"`
	Following int       `json:"following,omitzero"`
	CreatedAt time.Time `json:"created_at,omitzero"`
}
