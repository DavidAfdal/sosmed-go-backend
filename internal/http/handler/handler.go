package handler

import "github.com/davidafdal/post-app/pkg/validator"

type Handler struct {
	UserHandler    *UserHandler
	FeedHandler    *FeedHandler
	CommentHandler *CommentHandler
}

func NewHandler(userhHandler *UserHandler, feedHnadler *FeedHandler, commentHandler *CommentHandler) Handler {
	return Handler{
		UserHandler:    userhHandler,
		FeedHandler:    feedHnadler,
		CommentHandler: commentHandler,
	}
}

func checkValidation(input interface{}) (errorMessage string, data interface{}) {
	validationErrors := validator.Validate(input)
	if validationErrors != nil {
		if _, exists := validationErrors["error"]; exists {
			return "validasi input gagal", nil
		}
		return "validasi input gagal", validationErrors
	}
	return "", nil
}
