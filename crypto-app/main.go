package main

import (
	"log"

	"gses2.app/api/pkg"
	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/presentation/http/routes"
)

func main() {
	config.LoadEnv()
	log.Fatal(routes.SetupRouter(pkg.InitServices()).Run(config.NetworkPort))
}
