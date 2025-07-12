package server

import (
	"fmt"
	"log/slog"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/api"
	"github.com/aetheris-lab/aetheris-id/api/internal/handlers"
	"github.com/aetheris-lab/aetheris-id/api/internal/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type Server struct {
	echo *echo.Echo
	port string
}

func NewServer(config *configs.Environment, clientHandler handlers.ClientHandler, authHandler handlers.AuthHandler, oauthHandler handlers.OAuthHandler, authMiddleware middlewares.AuthMiddleware) *Server {
	e := echo.New()
	s := &Server{
		echo: e,
		port: fmt.Sprintf(":%d", config.Server.Port),
	}

	s.configureMiddlewares(config)
	s.configureValidator()
	s.configureErrorHandler()
	s.configureRoutes(config, clientHandler, authHandler, oauthHandler, authMiddleware)

	return s
}

func (s *Server) Start() error {
	return s.echo.Start(s.port)
}

func (s *Server) configureMiddlewares(config *configs.Environment) {
	s.echo.Use(middleware.Recover())

	s.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("request",
				"method", v.Method,
				"uri", v.URI,
				"status", v.Status,
				"error", v.Error,
			)
			return nil
		},
	}))

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.Cors.AllowedOrigins,
		AllowMethods:     config.Cors.AllowedMethods,
		AllowHeaders:     config.Cors.AllowedHeaders,
		AllowCredentials: true,
	}))

	s.echo.Use(middleware.BodyLimit("2M"))

	s.echo.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(config.RateLimit.MaxRequests),
	)))

}

func (s *Server) configureValidator() {
	s.echo.Validator = api.NewCustomValidator()
}

func (s *Server) configureErrorHandler() {
	s.echo.HTTPErrorHandler = api.CustomHTTPErrorHandler
}

func (s *Server) configureRoutes(config *configs.Environment, clientHandler handlers.ClientHandler, authHandler handlers.AuthHandler, oauthHandler handlers.OAuthHandler, authMiddleware middlewares.AuthMiddleware) {
	apiGroup := s.echo.Group("/api/v1")
	RegisterRoutes(apiGroup, config, clientHandler, authHandler, oauthHandler, authMiddleware)
}
