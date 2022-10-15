package main

import (
	"log"

	"customers-service/config"
	"customers-service/presentation"
)

func main() {
	router := presentation.GetRouter()

	log.Println("started")
	log.Fatal(router.Run(config.CustomersServerPort))
}
