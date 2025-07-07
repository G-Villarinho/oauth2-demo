package bootstrap

import (
	"github.com/aetheris-lab/aetheris-id/api/internal/handlers"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
	"github.com/aetheris-lab/aetheris-id/api/internal/server"
	"github.com/aetheris-lab/aetheris-id/api/internal/services"
	"github.com/aetheris-lab/aetheris-id/api/pkg/injector"
	"go.uber.org/dig"
)

func BuildContainer(container *dig.Container) {
	// Handlers
	injector.Provide(container, handlers.NewClientHandler)

	// Services
	injector.Provide(container, services.NewAuthService)
	injector.Provide(container, services.NewClientService)
	injector.Provide(container, services.NewJWTService)
	injector.Provide(container, services.NewOTPService)

	// Repositories
	injector.Provide(container, repositories.NewClientRepository)
	injector.Provide(container, repositories.NewOTPRepository)
	injector.Provide(container, repositories.NewUserRepository)

	// Server
	injector.Provide(container, server.NewServer)

}
