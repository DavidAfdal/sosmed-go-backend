package builder

import (
	"github.com/davidafdal/post-app/internal/http/handler"
	"github.com/davidafdal/post-app/internal/http/router"
	"github.com/davidafdal/post-app/internal/repositories"
	"github.com/davidafdal/post-app/internal/services"
	"github.com/davidafdal/post-app/pkg/cloudinary"
	"github.com/davidafdal/post-app/pkg/rabbitmq"
	"github.com/davidafdal/post-app/pkg/route"
	"github.com/davidafdal/post-app/pkg/token"
	"github.com/davidafdal/post-app/pkg/upload"
	"github.com/jmoiron/sqlx"
)

func BuildPublicRoute(db *sqlx.DB, cloudinary cloudinary.CloudinaryUseCase, token token.TokenUseCase) []*route.Route {

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cloudinary, token)
	userHandler := handler.NewUserHandler(userService)

	handler := handler.NewHandler(userHandler, nil, nil)

	return router.PublicRoute(handler)
}

func BuildPrivateRoute(db *sqlx.DB, cloudinary cloudinary.CloudinaryUseCase, token token.TokenUseCase, msgBroker *rabbitmq.Client) []*route.Route {
	uploadUsecase := upload.NewUploadUseCase()

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cloudinary, token)
	userHandler := handler.NewUserHandler(userService)

	commentRepo := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentService)

	feedRepo := repositories.NewFeedRepository(db)
	feedService := services.NewFeedService(feedRepo, uploadUsecase, msgBroker)
	feedHandler := handler.NewFeedHandler(feedService)

	handler := handler.NewHandler(userHandler, feedHandler, commentHandler)

	return router.PrivateRoute(handler)
}
