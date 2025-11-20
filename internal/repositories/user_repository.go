package repositories

import (
	"fmt"

	"github.com/davidafdal/post-app/internal/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Find(search string) ([]*entities.User, error)
	Create(user *entities.User) (*entities.User, error)
	FindByEmail(credentials string) (*entities.User, error)
	FindByUsername(username string) (*entities.User, error)
	FindByID(id uuid.UUID) (*entities.User, error)
	ToggleFollow(followerID, followingID uuid.UUID) (string, error)
	Update(user *entities.User) (*entities.User, error)
	Delete(userID uuid.UUID) error
}

type userRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepositoryImpl{db}
}

func (r *userRepositoryImpl) Find(search string) ([]*entities.User, error) {
	users := make([]*entities.User, 0)
	query := `
		SELECT id, username, email, avatar, created_at, updated_at
		FROM users
		WHERE username ILIKE $1
		ORDER BY username DESC
	`

	searchPattern := fmt.Sprintf("%%%s%%", search)

	if err := r.db.Select(&users, query, searchPattern); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepositoryImpl) Create(user *entities.User) (*entities.User, error) {
	newUser := new(entities.User)
	query := `
		INSERT INTO users (username, email, password, avatar)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, avatar, created_at, updated_at;
	`

	if err := r.db.Get(newUser, query, user.Username, user.Email, user.Password, user.Avatar); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *userRepositoryImpl) FindByEmail(credential string) (*entities.User, error) {
	user := new(entities.User)
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1 or username = $1;
	`

	if err := r.db.Get(user, query, credential); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByUsername(username string) (*entities.User, error) {
	user := new(entities.User)
	query := `
		SELECT 
			u.id,
			u.username,
			u.email,
			COALESCE(u.avatar, '') AS avatar,
    		COALESCE(u.bio, '') AS bio,
			u.created_at,
			COUNT(DISTINCT f1.follower_id) as followers_count,
			COUNT(DISTINCT f2.following_id) as followings_count
		FROM users u
		LEFT JOIN user_folows f1 on f1.following_id = u.id
		LEFT JOIN user_folows f2 on f2.follower_id = u.id
		WHERE u.username = $1
		GROUP BY u.id, u.username
	`
	if err := r.db.Get(user, query, username); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByID(id uuid.UUID) (*entities.User, error) {
	user := new(entities.User)
	query := `
		SELECT id, username, email, bio, avatar, created_at, updated_at
		FROM users
		WHERE id = $1;  
	`
	if err := r.db.Get(user, query, id); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) Update(user *entities.User) (*entities.User, error) {
	updatedUser := new(entities.User)
	query := `
		UPDATE users
		SET username = $1,
			avatar = $2,
			bio = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, username, email, avatar, bio, created_at, updated_at
	`
	if err := r.db.Get(updatedUser, query, user.Username, user.Avatar, user.Bio, user.ID); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (r *userRepositoryImpl) Delete(userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Queryx(query, userID)

	return err
}

func (r *userRepositoryImpl) ToggleFollow(followerID, followingID uuid.UUID) (string, error) {
	isFollowing, err := r.isFollowing(followerID, followingID)

	if err != nil {
		return "", err
	}

	if isFollowing {
		err = r.unfollowUser(followerID, followingID)
		if err != nil {
			return "", err
		}
		return "unfollowed", nil
	}

	err = r.followUser(followerID, followingID)

	if err != nil {
		return "", err
	}

	return "followed", nil
}

func (r *userRepositoryImpl) isFollowing(followerID, followingID uuid.UUID) (bool, error) {
	var exits bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM user_folows
			WHERE follower_id = $1 and following_id = $2
		);
	`
	err := r.db.Get(&exits, query, followerID, followingID)
	return exits, err
}

func (r *userRepositoryImpl) followUser(followerID, followingID uuid.UUID) error {
	query := `
		INSERT INTO user_folows (follower_id, following_id)
		VALUES ($1, $2)
		ON CONFLICT (follower_id, following_id) DO NOTHING;
	`
	_, err := r.db.Exec(query, followerID, followingID)
	return err
}

func (r *userRepositoryImpl) unfollowUser(followerID, followingID uuid.UUID) error {
	query := `
		DELETE FROM user_folows 
		WHERE follower_id = $1 AND following_id = $2;
	`
	_, err := r.db.Exec(query, followerID, followingID)
	return err
}
