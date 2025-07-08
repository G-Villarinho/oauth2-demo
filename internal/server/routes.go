package server

import (
	"github.com/aetheris-lab/aetheris-id/api/internal/handlers"
	"github.com/aetheris-lab/aetheris-id/api/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(apiGroup *echo.Group, clientHandler handlers.ClientHandler, authHandler handlers.AuthHandler, oauthHandler handlers.OAuthHandler, authMiddleware middlewares.AuthMiddleware) {
	registerClientRoutes(apiGroup, clientHandler)
	registerAuthRoutes(apiGroup, authHandler, authMiddleware)
	registerOAuthRoutes(apiGroup, oauthHandler, authMiddleware)
}

func registerClientRoutes(group *echo.Group, clientHandler handlers.ClientHandler) {
	group.POST("/clients", clientHandler.CreateClient)
}

func registerAuthRoutes(group *echo.Group, h handlers.AuthHandler, authMiddleware middlewares.AuthMiddleware) {
	group.POST("/auth/login", h.Login)
	group.POST("/auth/authenticate", h.Authenticate, authMiddleware.EnsureOTPAuthenticated())
}

func registerOAuthRoutes(group *echo.Group, h handlers.OAuthHandler, authMiddleware middlewares.AuthMiddleware) {
	oauthGroup := group.Group("/oauth")

	oauthGroup.GET("/authorize", h.Authorize, authMiddleware.AttachUserClaimsIfAuthenticated())
}
