package repositories

import (
	"context"
	"errors"

	"github.com/davidafdal/post-app/internal/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *entities.Comment) (*entities.Comment, error)
	CreateReply(ctx context.Context, comment *entities.Comment) (*entities.Comment, error)
	FindTopComment(ctx context.Context, feedID uuid.UUID) ([]*entities.Comment, error)
	FindRepliesComment(ctx context.Context, commentID uuid.UUID) ([]*entities.Comment, error)
}

type commentRepositoryImpl struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) CommentRepository {
	return &commentRepositoryImpl{
		db: db,
	}
}

func (r *commentRepositoryImpl) Create(ctx context.Context, comment *entities.Comment) (*entities.Comment, error) {
	query := `
		INSERT INTO feed_comments (feed_id, user_id, comment)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	err := r.db.QueryRowContext(ctx, query, comment.FeedID, comment.UserID, comment.Comment).Scan(&comment.ID)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *commentRepositoryImpl) CreateReply(ctx context.Context, comment *entities.Comment) (*entities.Comment, error) {
	if err := r.findByID(ctx, comment.ParentID); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO feed_comments (feed_id, user_id, parent_id, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	err := r.db.QueryRowContext(ctx, query, comment.FeedID, comment.UserID, comment.ParentID, comment.Comment).Scan(&comment.ID)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *commentRepositoryImpl) findByID(ctx context.Context, commentID uuid.UUID) error {
	var exits bool

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM feed_comments
			WHERE id = $1
		);
	`

	if err := r.db.QueryRowContext(ctx, query, commentID).Scan(&exits); err != nil {
		return err
	}

	if !exits {
		return errors.New("comment not found")
	}

	return nil
}

func (r *commentRepositoryImpl) FindTopComment(ctx context.Context, feedID uuid.UUID) ([]*entities.Comment, error) {
	comments := make([]*entities.Comment, 0)

	query := `
		SELECT 
			c.*,
			(
				SELECT COUNT(*)
				FROM feed_comments rc
				WHERE rc.parent_id = c.id
			) as reply_count
		FROM feed_comments c
		WHERE c.feed_id = $1 
			AND c.parent_id IS NULL
		ORDER BY c.created_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query, feedID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c entities.Comment

		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}

	return comments, nil
}

func (r *commentRepositoryImpl) FindRepliesComment(ctx context.Context, commentID uuid.UUID) ([]*entities.Comment, error) {
	comments := make([]*entities.Comment, 0)

	query := `
		SELECT * 
		FROM feed_comments
		WHERE parent_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, commentID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c entities.Comment

		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}

	return comments, nil
}
