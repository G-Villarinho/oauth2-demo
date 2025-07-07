package server

import (
	"github.com/aetheris-lab/aetheris-id/api/internal/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(apiGroup *echo.Group, clientHandler handlers.ClientHandler) {
	registerClientRoutes(apiGroup, clientHandler)
}

func registerClientRoutes(group *echo.Group, clientHandler handlers.ClientHandler) {
	group.POST("/clients", clientHandler.CreateClient)
}
