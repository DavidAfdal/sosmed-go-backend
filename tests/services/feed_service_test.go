package services_test

import (
	"context"
	"mime/multipart"
	"testing"

	"github.com/davidafdal/post-app/internal/dto"
	"github.com/davidafdal/post-app/internal/entities"
	"github.com/davidafdal/post-app/internal/events"
	"github.com/davidafdal/post-app/internal/services"
	mocksPkg "github.com/davidafdal/post-app/mocks/pkg"
	mocksRepo "github.com/davidafdal/post-app/mocks/repositories"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFeedService_CreateFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// mock dependencies
	feedRepo := mocksRepo.NewMockFeedRepository(ctrl)
	storage := mocksPkg.NewMockUploadUseCase(ctrl)
	publisher := mocksPkg.NewMockMessageBroker(ctrl)

	// service under test
	svc := services.NewFeedService(feedRepo, storage, publisher)

	req := &dto.CreateFeedRequest{
		Caption: "test caption",
		UserID:  uuid.New(),
	}

	// Dummy file header (valid & minimal)
	fileHeader := &multipart.FileHeader{
		Filename: "test.jpg",
		Size:     10,
	}
	files := []*multipart.FileHeader{fileHeader}

	mockFeed := &entities.Feed{
		ID:      uuid.New(),
		Caption: req.Caption,
		User: &entities.User{
			ID:       req.UserID,
			Username: "testuser",
			Avatar:   "avatar.jpg",
		},
		Medias: []*entities.FeedMedia{
			{
				Url:  "https://example.com/img.jpg",
				Type: "image",
			},
		},
		Likes:    5,
		Comments: 2,
	}

	// EXPECTATIONS
	feedRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(mockFeed, nil)

	storage.
		EXPECT().
		SaveTempFile(gomock.Any(), gomock.Any()).
		Return("http://example.com/image.jpg", nil)

	publisher.
		EXPECT().
		Publish(
			"",
			"events",
			string(events.UploadFeedMedias),
			gomock.Any(),
		).
		Return(nil)

	// DO
	res, err := svc.CreateFeed(ctx, req, files)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Caption, res.Caption)
}
