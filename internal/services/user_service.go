package services

import (
	"github.com/davidafdal/post-app/internal/entities"
	"github.com/davidafdal/post-app/internal/repositories"
)

type UserService interface {
}

type userServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) GetUsers(search string) ([]*entities.User, error) {
	users, err := s.userRepo.Find(search)

	if err != nil {
		return nil, err
	}

	return users, nil
}
