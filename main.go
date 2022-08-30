package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gses2.app/api/config"
	"gses2.app/api/handlers"
)

func main() {
	config.LoadEnv()
	handleRequests()
}

// handleRequests registers handlers and starts running the server.
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/rate", handlers.RateHandler).Methods("GET")
	router.HandleFunc("/subscribe", handlers.SubscribeHandler).Methods("POST")
	router.HandleFunc("/sendEmails", handlers.SendEmailsHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(config.Port, router))
}
