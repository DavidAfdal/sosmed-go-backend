package services

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/davidafdal/post-app/internal/dto"
	"github.com/davidafdal/post-app/internal/entities"
	"github.com/davidafdal/post-app/internal/events"
	"github.com/davidafdal/post-app/internal/repositories"
	"github.com/davidafdal/post-app/pkg/rabbitmq"
	"github.com/davidafdal/post-app/pkg/upload"
	"github.com/google/uuid"
)

type FeedService interface {
	CreateFeed(ctx context.Context, req *dto.CreateFeedRequest, files []*multipart.FileHeader) (*dto.FeedResponse, error)
	GetFeeds(ctx context.Context, userID uuid.UUID) ([]*dto.FeedResponse, error)
	LikeFeed(feedID, userID uuid.UUID) (string, error)
}

type feedServicesImpl struct {
	feedRepo      repositories.FeedRepository
	uploadUseCase upload.UploadUseCase
	msgBroker     rabbitmq.MessageBroker
}

func NewFeedService(feedRepo repositories.FeedRepository, uploadUseCase upload.UploadUseCase, msgBroker rabbitmq.MessageBroker) FeedService {
	return &feedServicesImpl{
		feedRepo:      feedRepo,
		uploadUseCase: uploadUseCase,
		msgBroker:     msgBroker,
	}
}

func (s *feedServicesImpl) CreateFeed(ctx context.Context, req *dto.CreateFeedRequest, files []*multipart.FileHeader) (*dto.FeedResponse, error) {

	feed := &entities.Feed{
		Caption: req.Caption,
		UserID:  req.UserID,
	}

	createdFeed, err := s.feedRepo.Create(ctx, feed)

	if err != nil {
		return nil, err
	}

	contentData := make([]events.ContentData, len(files))

	for i, fileHeader := range files {
		tempPath, err := s.uploadUseCase.SaveTempFile(fileHeader, "/app/uploads")
		if err != nil {
			return nil, err
		}
		contentData[i] = events.ContentData{
			FilePath: tempPath,
			FileType: "Image",
		}
	}

	payload := events.UploadPayload{
		FeedID:  createdFeed.ID.String(),
		Content: contentData,
	}

	body, _ := json.Marshal(payload)

	if err := s.msgBroker.Publish("", "events", string(events.UploadFeedMedias), body); err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			fmt.Println("delete data from image")
		}
	}()

	return s.toFeedResponse(createdFeed), nil
}

func (s *feedServicesImpl) GetFeeds(ctx context.Context, userID uuid.UUID) ([]*dto.FeedResponse, error) {
	feeds, err := s.feedRepo.GetFeeds(ctx, userID, 100)

	if err != nil {
		return nil, err
	}

	feedResponse := make([]*dto.FeedResponse, 0, len(feeds))

	for _, feed := range feeds {
		feedResponse = append(feedResponse, s.toFeedResponse(feed))
	}

	return feedResponse, nil
}

// func (s *feedServicesImpl) GetFeedByID(ctx context.Context, feedId uuid.UUID) (*dto.FeedResponse, error) {

// }

func (s *feedServicesImpl) LikeFeed(feedID, userID uuid.UUID) (string, error) {
	status, err := s.feedRepo.ToggleLiked(feedID, userID)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (s *feedServicesImpl) toFeedResponse(feed *entities.Feed) *dto.FeedResponse {
	mediasResponse := make([]*dto.MediaResponse, 0, len(feed.Medias))

	for _, media := range feed.Medias {
		mediasResponse = append(mediasResponse, &dto.MediaResponse{
			Url:  media.Url,
			Type: media.Type,
		})
	}

	userResponse := &dto.UserResponse{
		Username: feed.User.Username,
		Avatar:   feed.User.Avatar,
	}

	return &dto.FeedResponse{
		ID:       feed.ID,
		Caption:  feed.Caption,
		Medias:   mediasResponse,
		User:     userResponse,
		Likes:    feed.Likes,
		Comments: feed.Comments,
	}
}
