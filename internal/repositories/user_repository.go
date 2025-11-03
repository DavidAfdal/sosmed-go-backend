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
	FindByEmailOrUsername(credentials string) (*entities.User, error)
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
		VALUE ($1, $2, $3, $4)
		RETURNING id, username, email, avatar, created_at, updated_at
	`

	if err := r.db.Get(newUser, query, user.Username, user.Email, user.Password, user.Avatar); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *userRepositoryImpl) FindByEmailOrUsername(credential string) (*entities.User, error) {
	user := new(entities.User)
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1 or username = $1
	`

	if err := r.db.Get(user, query, credential); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindUserWithFeed(username string) (*entities.User, error) {

	return nil, nil
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
	query := `DELETE users WHERE id = $1`
	_, err := r.db.Queryx(query, userID)

	return err
}
