package handler

import (
	"net/http"

	"github.com/davidafdal/post-app/internal/dto"
	"github.com/davidafdal/post-app/internal/services"
	"github.com/davidafdal/post-app/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	searchQuery := c.QueryParam("search")
	users, err := h.userService.GetUsers(searchQuery)

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "succes get all users", users)
}

func (h *UserHandler) GetByUsername(c echo.Context) error {
	usernameParam := c.Param("username")

	user, err := h.userService.GetUserByUsername(usernameParam)

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "success get user by username", user)
}

func (h *UserHandler) Login(c echo.Context) error {

	req := new(dto.LoginRequest)

	if err := c.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if errMessage, data := checkValidation(req); errMessage != "" {
		return response.SuccessResponse(c, http.StatusBadRequest, errMessage, data)
	}

	responData, err := h.userService.Login(req)

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "success login by credentials", responData)
}

func (h *UserHandler) Register(c echo.Context) error {
	req := new(dto.CreateUserRequest)

	if err := c.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if errMessage, data := checkValidation(req); errMessage != "" {
		return response.SuccessResponse(c, http.StatusBadRequest, errMessage, data)
	}

	avatarFile, err := c.FormFile("avatar")

	if err != nil && err != http.ErrMissingFile {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	user, err := h.userService.Register(req, avatarFile)

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, "success registed user", user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Get("user_id").(string)
	userID := uuid.MustParse(id)
	req := new(dto.UpdatedUserRequest)

	if err := c.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if errMessage, data := checkValidation(*req); errMessage != "" {
		return response.SuccessResponse(c, http.StatusBadRequest, errMessage, data)
	}

	avatarFile, err := c.FormFile("avatar")

	if err != nil && err != http.ErrMissingFile {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	user, err := h.userService.UpdateUser(req, avatarFile, userID)

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "success update data user", user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Get("user_id").(string)
	userID := uuid.MustParse(id)

	if err := h.userService.DeleteUser(userID); err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "sucess delete user", nil)
}

func (h *FeedHandler) LikeFeed(c echo.Context) error {
	id := c.Get("user_id").(string)
	paramId := c.Param("feed_id")

	userID := uuid.MustParse(id)
	feedID := uuid.MustParse(paramId)

	status, err := h.feedService.LikeFeed(feedID, userID)

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "succes likes feed", map[string]interface{}{
		"status": status,
	})
}
