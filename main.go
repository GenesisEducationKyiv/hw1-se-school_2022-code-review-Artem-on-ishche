package main

import (
	"gses2.app/api/config"
	"gses2.app/api/routes"
)

func main() {
	config.LoadEnv()
	routes.RegisterRoutes()
}
