package main

import (
	"log"

	"gses2.app/api/pkg"
	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/infrastructure/logger"
	"gses2.app/api/pkg/presentation/http/routes"
)

func main() {
	loggerService := logger.NewRabbitMQLogger()
	defer loggerService.Close()

	config.LoadEnv(loggerService)
	log.Fatal(routes.SetupRouter(pkg.InitServices(loggerService)).Run(config.NetworkPort))
}
