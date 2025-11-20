package handler

import (
	"net/http"

	"github.com/davidafdal/post-app/internal/dto"
	"github.com/davidafdal/post-app/internal/services"
	"github.com/davidafdal/post-app/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	commentService services.CommentService
}

func NewCommentHandler(commentService services.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) CreateComment(c echo.Context) error {
	payloadID := c.Get("user_id").(string)
	userID := uuid.MustParse(payloadID)
	req := new(dto.CreateCommentRequest)

	if err := c.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if errMsg, data := checkValidation(req); errMsg != "" {
		return response.SuccessResponse(c, http.StatusBadRequest, errMsg, data)
	}

	req.SenderID = userID

	if err := h.commentService.CreateComment(c.Request().Context(), req); err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, "success create a comment", nil)
}

func (h *CommentHandler) CreateReplyComment(c echo.Context) error {
	payloadID := c.Get("user_id").(string)
	senderID := uuid.MustParse(payloadID)

	req := new(dto.CreateReplyCommentRequest)

	if err := c.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if errMsg, data := checkValidation(req); errMsg != "" {
		return response.SuccessResponse(c, http.StatusBadRequest, errMsg, data)
	}

	req.SenderID = senderID

	if err := h.commentService.CreateCommentReplies(c.Request().Context(), req); err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, "succes create a new comment", nil)
}

func (h *CommentHandler) GetTopLevelComment(c echo.Context) error {
	feedID := uuid.MustParse(c.Param("feed_id"))

	responData, err := h.commentService.GetTopLevelComment(c.Request().Context(), feedID)

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, "success get top level comment", responData)
}

func (h *CommentHandler) GetCommentReplies(c echo.Context) error {
	commentID := uuid.MustParse(c.Param("comment_id"))

	responData, err := h.commentService.GetRepliedComment(c.Request().Context(), commentID)

	if err != nil {
		return response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, "success get top reply comment", responData)
}
