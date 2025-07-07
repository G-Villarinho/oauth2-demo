package api

import (
	"github.com/aetheris-lab/aetheris-id/api/internal/handlers"
	"github.com/aetheris-lab/aetheris-id/api/internal/repositories"
	"github.com/aetheris-lab/aetheris-id/api/internal/services"
	"github.com/aetheris-lab/aetheris-id/api/pkg/injector"
	"go.uber.org/dig"
)

func InitializeInternalDependencies(container *dig.Container) {
	provideHandlers(container)
	provideServices(container)
	provideRepositories(container)
}

func provideHandlers(container *dig.Container) {
	injector.Provide(container, handlers.NewClientHandler)
}

func provideServices(container *dig.Container) {
	injector.Provide(container, services.NewAuthService)
	injector.Provide(container, services.NewClientService)
	injector.Provide(container, services.NewJWTService)
	injector.Provide(container, services.NewOTPService)
}

func provideRepositories(container *dig.Container) {
	injector.Provide(container, repositories.NewClientRepository)
	injector.Provide(container, repositories.NewOTPRepository)
	injector.Provide(container, repositories.NewUserRepository)
}
