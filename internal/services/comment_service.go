package services

import (
	"context"

	"github.com/davidafdal/post-app/internal/dto"
	"github.com/davidafdal/post-app/internal/entities"
	"github.com/davidafdal/post-app/internal/repositories"
	"github.com/google/uuid"
)

type CommentService interface {
	CreateCommentReplies(ctx context.Context, req *dto.CreateReplyCommentRequest) error
	GetTopLevelComment(ctx context.Context, feedID uuid.UUID) ([]*dto.CommentResponse, error)
	GetRepliedComment(ctx context.Context, commentID uuid.UUID) ([]*dto.CommentResponse, error)
	CreateComment(ctx context.Context, req *dto.CreateCommentRequest) error
}

type commentServiceImpl struct {
	commentRepo repositories.CommentRepository
}

func NewCommentService(commentRepo repositories.CommentRepository) CommentService {
	return &commentServiceImpl{
		commentRepo: commentRepo,
	}
}

func (s *commentServiceImpl) CreateComment(ctx context.Context, req *dto.CreateCommentRequest) error {

	comment := &entities.Comment{
		FeedID:  req.FeedID,
		UserID:  req.SenderID,
		Comment: req.Comment,
	}

	_, err := s.commentRepo.Create(ctx, comment)

	if err != nil {
		return err
	}

	return nil
}

func (s *commentServiceImpl) CreateCommentReplies(ctx context.Context, req *dto.CreateReplyCommentRequest) error {
	commentData := &entities.Comment{
		FeedID:   req.FeedID,
		ParentID: req.CommentID,
		UserID:   req.SenderID,
		Comment:  req.Comment,
	}
	_, err := s.commentRepo.CreateReply(ctx, commentData)

	if err != nil {
		return err
	}

	return nil
}

func (s *commentServiceImpl) GetTopLevelComment(ctx context.Context, feedID uuid.UUID) ([]*dto.CommentResponse, error) {
	comments, err := s.commentRepo.FindTopComment(ctx, feedID)

	if err != nil {
		return nil, err
	}

	commentsResponse := make([]*dto.CommentResponse, len(comments))

	for i, v := range comments {
		commentsResponse[i] = s.toCommentResponse(v)
	}

	return commentsResponse, nil
}

func (s *commentServiceImpl) GetRepliedComment(ctx context.Context, commentID uuid.UUID) ([]*dto.CommentResponse, error) {
	replycomments, err := s.commentRepo.FindRepliesComment(ctx, commentID)

	if err != nil {
		return nil, err
	}

	responseComment := make([]*dto.CommentResponse, len(replycomments))

	for i, v := range replycomments {
		responseComment[i] = s.toCommentResponse(v)
	}

	return responseComment, nil
}

func (r *commentServiceImpl) toCommentResponse(comment *entities.Comment) *dto.CommentResponse {

	return &dto.CommentResponse{
		ID: comment.ID,
	}
}
