package main

import (
	"log"

	"github.com/aetheris-lab/aetheris-id/api/configs"
	"github.com/aetheris-lab/aetheris-id/api/infra/database"
	"github.com/aetheris-lab/aetheris-id/api/internal/bootstrap"
	"github.com/aetheris-lab/aetheris-id/api/internal/server"
	"github.com/aetheris-lab/aetheris-id/api/pkg/injector"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()

	injector.Provide(container, configs.NewConfig)

	injector.Provide(container, database.NewMongoClient)
	injector.Provide(container, database.NewMongoDatabase)

	bootstrap.BuildContainer(container)

	server := injector.Resolve[*server.Server](container)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
