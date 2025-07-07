package server

import (
	"fmt"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
	port string
}

func NewServer(config *configs.Environment, clientHandler handlers.ClientHandler) *Server {
	e := echo.New()

	e.Use(middleware.Recover())

	return &Server{
		echo: e,
		port: fmt.Sprintf(":%d", config.Server.Port),
	}
}
