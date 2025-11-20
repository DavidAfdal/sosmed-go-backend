package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/davidafdal/post-app/pkg/response"
	"github.com/davidafdal/post-app/pkg/route"
	"github.com/davidafdal/post-app/pkg/token"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	*echo.Echo
}

func NewServer(publicRoutes, privateRoutes []*route.Route, secretKey string, tokenUse token.TokenUseCase) *Server {
	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return response.SuccessResponse(c, http.StatusOK, "Hello, World!", nil)
	})

	v1 := e.Group("api")

	if len(publicRoutes) > 0 {
		for _, v := range publicRoutes {
			v1.Add(v.Method, v.Path, v.Handler)
		}
	}

	if len(privateRoutes) > 0 {
		for _, v := range privateRoutes {
			v1.Add(v.Method, v.Path, v.Handler, JWTProtection(secretKey), UserContextMiddelware())
		}
	}

	return &Server{e}
}

func (s *Server) Run() {
	runServer(s)
	gracefulShutdown(s)
}

func runServer(srv *Server) {
	go func() {
		err := srv.Start(":8080")
		log.Fatal(err)
	}()
}

func gracefulShutdown(srv *Server) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		if err := srv.Shutdown(ctx); err != nil {
			srv.Logger.Fatal("Server Shutdown:", err)
		}
	}()
}

func UserContextMiddelware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			user, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return response.ErrorResponse(c, http.StatusUnauthorized, "anda harus login untuk mengakses resource ini")
			}

			claims := user.Claims.(*token.JwtCustomClaims)

			c.Set("user_id", claims.ID)
			c.Set("user_email", claims.Email)

			return next(c)
		}
	}
}

func JWTProtection(secretKey string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		SigningKey: []byte(secretKey),
		ErrorHandler: func(c echo.Context, err error) error {
			return response.ErrorResponse(c, http.StatusUnauthorized, "anda harus login untuk mengakses resource ini")
		},
	})
}
