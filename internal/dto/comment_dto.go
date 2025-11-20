package dto

import "github.com/google/uuid"

type CreateCommentRequest struct {
	FeedID   uuid.UUID `param:"feed_id" validate:"required"`
	SenderID uuid.UUID
	Comment  string `json:"comment" validate:"required"`
}

type CreateReplyCommentRequest struct {
	CommentID uuid.UUID `param:"comment_id" validate:"required"`
	FeedID    uuid.UUID `json:"feed_id" validate:"requried"`
	SenderID  uuid.UUID
	Comment   string `json:"comment" validate:"required"`
}

type CommentResponse struct {
	ID         uuid.UUID     `json:"id"`
	Comment    string        `json:"comment"`
	User       *UserResponse `json:"user"`
	ReplyCount int           `json:"replies,omizero"`
}
