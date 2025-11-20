package main

import (
	"time"

	"github.com/davidafdal/post-app/config"
	"github.com/davidafdal/post-app/internal/builder"
	"github.com/davidafdal/post-app/pkg/cloudinary"
	"github.com/davidafdal/post-app/pkg/postgres"
	"github.com/davidafdal/post-app/pkg/rabbitmq"
	"github.com/davidafdal/post-app/pkg/server"
	"github.com/davidafdal/post-app/pkg/token"
)

func main() {
	cfg, err := config.NewConfig()
	checkError(err)
	db, err := postgres.InitPostgres(&cfg.Postgres)
	checkError(err)
	clodinary, err := cloudinary.NewCloudinaryUseCase(&cfg.Cloudinary)
	checkError(err)
	token := token.NewTokenUseCase(cfg.JWT.SecretKey, time.Duration(cfg.JWT.ExpiresAt)*time.Hour)

	rqm, err := rabbitmq.NewClient(&cfg.Rabbit)
	checkError(err)

	publicRoutes := builder.BuildPublicRoute(db, clodinary, token)
	privateRoutes := builder.BuildPrivateRoute(db, clodinary, token, rqm)

	srv := server.NewServer(publicRoutes, privateRoutes, cfg.JWT.SecretKey, token)
	srv.Run()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
