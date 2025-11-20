package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/davidafdal/post-app/internal/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type feedRow struct {
	FeedID    uuid.UUID `db:"id"`
	Caption   string    `db:"caption"`
	CreatedAt time.Time `db:"created_at"`

	UserID   uuid.UUID `db:"user_id"`
	Username string    `db:"username"`
	Avatar   string    `db:"avatar"`

	MediaURL  string `db:"url"`
	MediaType string `db:"type"`

	Likes    int `db:"likes"`
	Comments int `db:"comments"`
}

type FeedRepository interface {
	Create(ctx context.Context, feed *entities.Feed) (*entities.Feed, error)
	GetFeeds(ctx context.Context, userID uuid.UUID, limit int) ([]*entities.Feed, error)
	ToggleLiked(feedID, userID uuid.UUID) (string, error)
}

type feedRepositoryImpl struct {
	db *sqlx.DB
}

func NewFeedRepository(db *sqlx.DB) FeedRepository {
	return &feedRepositoryImpl{db: db}
}

func (r *feedRepositoryImpl) Create(ctx context.Context, feed *entities.Feed) (*entities.Feed, error) {

	tx, err := r.db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var feedId uuid.UUID

	query := `
		INSERT INTO feeds (user_id, caption)
		VALUES ($1, $2)
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, query, feed.UserID.String(), feed.Caption).Scan(&feedId)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	feed.ID = feedId

	return feed, nil
}

func (r *feedRepositoryImpl) GetFeeds(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) ([]*entities.Feed, error) {

	query := `
		SELECT 
			f.id,
			f.caption,
			f.created_at,
			u.username,
			u.avatar,
			fm.url,
			fm.type,
			(
			  SELECT COUNT(*) 
			  FROM feed_likes fl 
			  WHERE fl.feed_id = f.id
			) AS likes,
			(
			  SELECT COUNT(*) 
			  FROM feed_comments fc 
			  WHERE fc.feed_id = f.id
			) AS comments
		FROM feeds f
		JOIN users u ON u.id = f.user_id
		JOIN feed_media fm ON fm.feed_id = f.id
		WHERE 
			(f.user_id IN (
				SELECT following_id FROM user_folows WHERE follower_id = $1
			) OR f.user_id = $1)
		ORDER BY f.created_at DESC
		LIMIT $2;
	`

	rows := make([]feedRow, 0)

	if err := r.db.SelectContext(ctx, &rows, query, userID, limit); err != nil {
		return nil, err
	}

	feedMap := make(map[uuid.UUID]*entities.Feed)

	for _, row := range rows {
		fmt.Println("contoh feedmap", feedMap[row.FeedID])

		if _, ok := feedMap[row.FeedID]; !ok {
			feedMap[row.FeedID] = &entities.Feed{
				ID:      row.FeedID,
				Caption: row.Caption,
				User: &entities.User{
					ID:       row.UserID,
					Username: row.Username,
					Avatar:   row.Avatar,
				},
				Likes:    row.Likes,
				Comments: row.Comments,
				Medias:   []*entities.FeedMedia{},
			}
		}

		feedMap[row.FeedID].Medias = append(feedMap[row.FeedID].Medias, &entities.FeedMedia{
			Url:  row.MediaURL,
			Type: row.MediaType,
		})
	}

	result := make([]*entities.Feed, 0, len(feedMap))

	for _, f := range feedMap {
		result = append(result, f)
	}

	return result, nil
}

func (r *feedRepositoryImpl) GetFeed(ctx context.Context, feedID uuid.UUID) (*entities.Feed, error) {
	query := `
		SELECT 
			f.id,
			f.caption,
			f.created_at,
			u.username,
			u.avatar,
			fm.url,
			fm.type,
			c.comment,
			(SELECT COUNT (*)
			 FROM feed_likes f1
			 WHERE f1.feed_id = f.id
			) as likes,
		FROM f
		JOIN users u ON u.id = f.user_id
		LEFT JOIN feed_medias fm ON fm.feed_id = f.id
		LEFT JOIN feed_comments c ON c.feed_id = f.id
		WHERE f.id = $1	
	`
	var row feedRow

	err := r.db.QueryRowContext(ctx, query, feedID).Scan(&row)

	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (r *feedRepositoryImpl) ToggleLiked(feedID, userID uuid.UUID) (string, error) {
	isLiked, err := r.isLiked(feedID, userID)

	if err != nil {
		return "", err
	}

	if isLiked {
		err = r.unlikeFeed(feedID, userID)
		if err != nil {
			return "", err
		}
		return "unliked", nil
	}

	err = r.likeFeed(feedID, userID)

	if err != nil {
		return "", err
	}

	return "liked", nil
}

func (r *feedRepositoryImpl) isLiked(feedID, userID uuid.UUID) (bool, error) {
	var exits bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM feed_likes
			WHERE feed_id = $1 and user_id = $2
		);
	`
	err := r.db.Get(&exits, query, feedID, userID)
	return exits, err
}

func (r *feedRepositoryImpl) likeFeed(feedID, userID uuid.UUID) error {
	query := `
		INSERT INTO feed_likes (feed_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (feed_id, user_id) DO NOTHING;
	`
	_, err := r.db.Exec(query, feedID, userID)
	return err
}

func (r *feedRepositoryImpl) unlikeFeed(feedID, userID uuid.UUID) error {
	query := `
		DELETE FROM feed_likes 
		WHERE feed_id = $1 AND user_id = $2;
	`
	_, err := r.db.Exec(query, feedID, userID)
	return err
}
