package main

import (
	"gses2.app/api/pkg/presentation/http/routes"
	"log"

	"gses2.app/api/pkg"
	"gses2.app/api/pkg/config"
)

func main() {
	config.LoadEnv()
	log.Fatal(routes.SetupRouter(pkg.InitServices()).Run(config.NetworkPort))
}
