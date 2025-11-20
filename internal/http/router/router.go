package router

import (
	"net/http"

	"github.com/davidafdal/post-app/internal/http/handler"
	"github.com/davidafdal/post-app/pkg/route"
)

func PublicRoute(handler handler.Handler) []*route.Route {
	userHandler := handler.UserHandler

	return []*route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: userHandler.GetAllUsers,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: userHandler.Register,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:username",
			Handler: userHandler.GetByUsername,
		},
	}
}

func PrivateRoute(handler handler.Handler) []*route.Route {
	userHandler := handler.UserHandler
	feedHandler := handler.FeedHandler
	commentHandler := handler.CommentHandler

	return []*route.Route{
		{
			Method:  http.MethodPut,
			Path:    "/users",
			Handler: userHandler.UpdateUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users",
			Handler: userHandler.DeleteUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/users/:following_id/follow",
			Handler: userHandler.FollowingUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/feeds",
			Handler: feedHandler.CreateFeed,
		},
		{
			Method:  http.MethodGet,
			Path:    "/feeds",
			Handler: feedHandler.GetFeeds,
		},
		{
			Method:  http.MethodPost,
			Path:    "/feeds/:feed_id/like",
			Handler: feedHandler.LikeFeed,
		},
		{
			Method:  http.MethodPost,
			Path:    "/feeds/:feed_id/comment",
			Handler: commentHandler.CreateComment,
		},
		{
			Method:  http.MethodGet,
			Path:    "/feeds/:feed_id/comment",
			Handler: commentHandler.GetTopLevelComment,
		},
		{
			Method:  http.MethodGet,
			Path:    "/comments/:comment_id/reply",
			Handler: commentHandler.GetCommentReplies,
		},
		{
			Method:  http.MethodPost,
			Path:    "/comments/:comment_id/reply",
			Handler: commentHandler.CreateReplyComment,
		},
	}
}
