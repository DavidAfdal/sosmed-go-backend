package events

import "encoding/json"

type EventType string

const (
	UploadFeedMedias EventType = "upload_feed_medias"
	DeleteFeedMedias EventType = "delete_feed_medias"
)

type RawEvent struct {
	EventType EventType       `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

type UploadPayload struct {
	FeedID  string        `json:"feed_id"`
	Content []ContentData `json:"content"`
}

type ContentData struct {
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
}

type DeletePayload struct {
	PublicID string `json:"public_id"`
}
