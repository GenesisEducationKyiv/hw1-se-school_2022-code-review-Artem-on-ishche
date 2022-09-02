package main

import (
	"gses2.app/api/config"
	"gses2.app/api/handlers"
)

func main() {
	config.LoadEnv()
	handlers.HandleRequests()
}
