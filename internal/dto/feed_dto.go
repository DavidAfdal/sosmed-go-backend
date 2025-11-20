package dto

import "github.com/google/uuid"

type CreateFeedRequest struct {
	Caption string `form:"caption" validate:"required"`
	UserID  uuid.UUID
}

type FeedResponse struct {
	ID       uuid.UUID        `json:"id"`
	Caption  string           `json:"caption"`
	User     *UserResponse    `json:"user,omitzero"`
	Medias   []*MediaResponse `json:"medias,omitzero"`
	Likes    int              `json:"likes"`
	Comments int              `json:"comments"`
}

type MediaResponse struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}
