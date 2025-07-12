package bootstrap

import (
	"github.com/aetheris-lab/aetheris-id/api/internal/handlers"
	"github.com/aetheris-lab/aetheris-id/api/internal/middlewares"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
	"github.com/aetheris-lab/aetheris-id/api/internal/server"
	"github.com/aetheris-lab/aetheris-id/api/internal/services"
	"github.com/aetheris-lab/aetheris-id/api/pkg/ecdsa"
	"github.com/aetheris-lab/aetheris-id/api/pkg/injector"
	"go.uber.org/dig"
)

func BuildContainer(container *dig.Container) {

	// Crypto
	injector.Provide(container, ecdsa.NewEcdsaKeyPair)

	// Handlers
	injector.Provide(container, handlers.NewAuthHandler)
	injector.Provide(container, handlers.NewClientHandler)
	injector.Provide(container, handlers.NewOAuthHandler)

	// Services
	injector.Provide(container, services.NewAuthService)
	injector.Provide(container, services.NewAuthorizationCodeService)
	injector.Provide(container, services.NewClientService)
	injector.Provide(container, services.NewJWTService)
	injector.Provide(container, services.NewOAuthService)
	injector.Provide(container, services.NewOTPService)
	injector.Provide(container, services.NewRefreshTokenService)

	// Repositories
	injector.Provide(container, repositories.NewAuthorizationCodeRepository)
	injector.Provide(container, repositories.NewClientRepository)
	injector.Provide(container, repositories.NewOTPRepository)
	injector.Provide(container, repositories.NewRefreshTokenRepository)
	injector.Provide(container, repositories.NewUserRepository)

	// Server
	injector.Provide(container, server.NewServer)

	// Middlewares
	injector.Provide(container, middlewares.NewAuthMiddleware)
	injector.Provide(container, middlewares.NewCookieMiddleware)
}
