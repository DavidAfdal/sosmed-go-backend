package handler

import (
	"net/http"

	"github.com/davidafdal/post-app/internal/dto"
	"github.com/davidafdal/post-app/internal/services"
	"github.com/davidafdal/post-app/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type FeedHandler struct {
	feedService services.FeedService
}

func NewFeedHandler(feedService services.FeedService) *FeedHandler {
	return &FeedHandler{
		feedService: feedService,
	}
}

func (h *FeedHandler) CreateFeed(c echo.Context) error {
	id := c.Get("user_id").(string)
	userID := uuid.MustParse(id)
	req := new(dto.CreateFeedRequest)

	if err := c.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if errMessage, data := checkValidation(req); errMessage != "" {
		return response.SuccessResponse(c, http.StatusBadRequest, errMessage, data)
	}

	req.UserID = userID

	form, err := c.MultipartForm()

	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, "invalid multipart form")
	}

	files := form.File["files"]

	if len(files) == 0 {
		return response.ErrorResponse(c, http.StatusBadRequest, "no files uploaded")
	}

	feed, err := h.feedService.CreateFeed(c.Request().Context(), req, files)

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "success create feed", feed)
}

func (h *FeedHandler) GetFeeds(c echo.Context) error {
	id := c.Get("user_id").(string)
	userID := uuid.MustParse(id)

	feeds, err := h.feedService.GetFeeds(c.Request().Context(), userID)

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "success get home feeds", feeds)
}

func (h *UserHandler) FollowingUser(c echo.Context) error {
	id := c.Get("user_id").(string)
	paramId := c.Param("following_id")

	followerID := uuid.MustParse(id)
	followingID := uuid.MustParse(paramId)

	status, err := h.userService.FollowUser(followerID, followingID)

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "succes following user", map[string]interface{}{
		"status": status,
	})
}
