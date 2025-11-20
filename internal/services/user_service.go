package services

import (
	"context"
	"mime/multipart"

	"github.com/davidafdal/post-app/internal/dto"
	"github.com/davidafdal/post-app/internal/entities"
	"github.com/davidafdal/post-app/internal/repositories"
	"github.com/davidafdal/post-app/pkg/cloudinary"
	"github.com/davidafdal/post-app/pkg/token"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type jwtResponse struct {
	Token      string `json:"token"`
	Expired_at string `json:"expired_at"`
}

type UserService interface {
	GetUsers(search string) ([]*dto.UserResponse, error)
	GetUserByUsername(username string) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (*jwtResponse, error)
	Register(req *dto.CreateUserRequest, file *multipart.FileHeader) (*dto.UserResponse, error)
	UpdateUser(req *dto.UpdatedUserRequest, file *multipart.FileHeader, userID uuid.UUID) (*dto.UserResponse, error)
	FollowUser(followerID, followingID uuid.UUID) (string, error)
	DeleteUser(userID uuid.UUID) error
}

type userServiceImpl struct {
	userRepo          repositories.UserRepository
	cloudinaryUseCase cloudinary.CloudinaryUseCase
	tokenUseCase      token.TokenUseCase
}

func NewUserService(userRepo repositories.UserRepository, cloudinaryUseCase cloudinary.CloudinaryUseCase, token token.TokenUseCase) UserService {
	return &userServiceImpl{
		userRepo:          userRepo,
		cloudinaryUseCase: cloudinaryUseCase,
		tokenUseCase:      token,
	}
}

func (s *userServiceImpl) GetUsers(search string) ([]*dto.UserResponse, error) {
	users, err := s.userRepo.Find(search)

	if err != nil {
		return nil, err
	}
	usersResponse := make([]*dto.UserResponse, len(users))

	for index, user := range users {
		usersResponse[index] = s.toUserResponse(user)
	}

	return usersResponse, nil
}

func (s *userServiceImpl) GetUserByUsername(username string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByUsername(username)

	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), err
}

func (s *userServiceImpl) Register(req *dto.CreateUserRequest, file *multipart.FileHeader) (*dto.UserResponse, error) {

	req.Avatar = ""

	if file != nil {
		filePath, err := s.cloudinaryUseCase.UploadFile(context.Background(), file)
		if err != nil {
			return nil, err
		}
		req.Avatar = filePath
	}

	hashPassowrd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)

	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashPassowrd),
		Avatar:   req.Avatar,
	}

	registerdUser, err := s.userRepo.Create(user)

	if err != nil {
		if req.Avatar != "" {
			_ = s.cloudinaryUseCase.DeleteFile(context.Background(), req.Avatar)
		}
		return nil, err
	}

	return s.toUserResponse(registerdUser), nil
}

func (s *userServiceImpl) Login(req *dto.LoginRequest) (*jwtResponse, error) {
	existedUser, err := s.userRepo.FindByEmail(req.Email)

	data := jwtResponse{}
	data.Token = ""
	data.Expired_at = ""

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	claims := s.tokenUseCase.CreateClaims(existedUser.ID.String(), existedUser.Email)

	accessToken, expiredAt, err := s.tokenUseCase.GenerateAccessToken(claims)

	if err != nil {
		return nil, err
	}

	data.Token = accessToken
	data.Expired_at = expiredAt.String()

	return &data, nil
}

func (s *userServiceImpl) UpdateUser(req *dto.UpdatedUserRequest, file *multipart.FileHeader, userID uuid.UUID) (*dto.UserResponse, error) {
	exits, err := s.userRepo.FindByID(userID)

	if err != nil {
		return nil, err
	}

	if file != nil {
		if exits.Avatar != "" {
			if err := s.cloudinaryUseCase.DeleteFile(context.Background(), exits.Avatar); err != nil {
				return nil, err
			}
		}

		secureUrl, err := s.cloudinaryUseCase.UploadFile(context.Background(), file)

		if err != nil {
			return nil, err
		}
		exits.Avatar = secureUrl
	}

	if req.Username != "" {
		exits.Username = req.Username
	}

	if req.Bio != "" {
		exits.Bio = req.Bio
	}

	user, err := s.userRepo.Update(exits)

	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *userServiceImpl) DeleteUser(userID uuid.UUID) error {
	exits, err := s.userRepo.FindByID(userID)

	if err != nil {
		return err
	}

	if exits.Avatar != "" {
		if err := s.cloudinaryUseCase.DeleteFile(context.Background(), exits.Avatar); err != nil {
			return err
		}
	}

	return s.userRepo.Delete(userID)
}

func (s *userServiceImpl) FollowUser(followerID, followingID uuid.UUID) (string, error) {
	status, err := s.userRepo.ToggleFollow(followerID, followingID)

	if err != nil {
		return "", err
	}

	return status, nil
}

func (s *userServiceImpl) toUserResponse(user *entities.User) *dto.UserResponse {
	dataResponse := &dto.UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Followers: user.Followers,
		Following: user.Followings,
		CreatedAt: user.CreatedAt,
	}

	return dataResponse
}
