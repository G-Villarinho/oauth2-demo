package server

import (
	"net/http"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/internal/handlers"
	"github.com/aetheris-lab/aetheris-id/api/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(apiGroup *echo.Group, env *configs.Environment, clientHandler handlers.ClientHandler, authHandler handlers.AuthHandler, oauthHandler handlers.OAuthHandler, authMiddleware middlewares.AuthMiddleware) {
	registerClientRoutes(apiGroup, clientHandler)
	registerAuthRoutes(apiGroup, authHandler, authMiddleware)
	registerOAuthRoutes(apiGroup, oauthHandler, authMiddleware)
	registerDevRoutes(apiGroup, env)
}

func registerDevRoutes(group *echo.Group, env *configs.Environment) {
	if env.Env == "development" {
		group.GET("/dev/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		})

		group.GET("/dev/envs", func(c echo.Context) error {
			return c.JSON(http.StatusOK, env)
		})
	}
}

func registerClientRoutes(group *echo.Group, clientHandler handlers.ClientHandler) {
	group.POST("/clients", clientHandler.CreateClient)
}

func registerAuthRoutes(group *echo.Group, h handlers.AuthHandler, authMiddleware middlewares.AuthMiddleware) {
	authGroup := group.Group("/auth")

	authGroup.POST("/login", h.Login)
	authGroup.POST("/register", h.Register)
	authGroup.POST("/authenticate", h.Authenticate, authMiddleware.EnsureOTPAuthenticated())
	authGroup.POST("/code/resend", h.ResendVerificationCode, authMiddleware.EnsureOTPAuthenticated())
}

func registerOAuthRoutes(group *echo.Group, h handlers.OAuthHandler, authMiddleware middlewares.AuthMiddleware) {
	oauthGroup := group.Group("/oauth")

	oauthGroup.GET("/authorize", h.Authorize, authMiddleware.AttachUserClaimsIfAuthenticated())
	oauthGroup.POST("/token", h.Token)
}
