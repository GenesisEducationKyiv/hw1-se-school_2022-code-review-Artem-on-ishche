package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gses2.app/api/config"
)

func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/rate", rateHandler).Methods("GET")
	router.HandleFunc("/subscribe", subscribeHandler).Methods("POST")
	router.HandleFunc("/sendEmails", sendEmailsHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(config.Port, router))
}
