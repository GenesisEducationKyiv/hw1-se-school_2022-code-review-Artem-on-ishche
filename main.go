package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gses2.app/api/data"
	"gses2.app/api/emails"
	"gses2.app/api/exchange_rate"
	"log"
	"net/http"
	"os"
)

func main() {
	loadEnv()
	handleRequests()
}

// loadEnv loads file .env to get environment variables from it later.
func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// handleRequests registers handlers and starts running the server.
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/rate", rateHandler).Methods("GET")
	router.HandleFunc("/subscribe", subscribeHandler).Methods("POST")
	router.HandleFunc("/sendEmails", sendEmailsHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	log.Fatal(http.ListenAndServe(port, router))
}

// rateHandler is responsible for GET requests at /rate.
func rateHandler(w http.ResponseWriter, _ *http.Request) {
	rate, err := exchange_rate.GetBtcUahRate()
	if err != nil || rate == -1 {
		http.Error(w, "Error has occurred", http.StatusBadRequest)
		return
	}
	_, _ = fmt.Fprint(w, rate)
}

// subscribeHandler is responsible for POST requests at /subscribe
func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	// get the parameter from request
	emailParams, ok := r.URL.Query()["email"]

	// if the request doesn't have a required parameter
	if !ok || len(emailParams[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "Bad Request")
		return
	}

	// save only the first email if more are provided
	email := emailParams[0]

	// check if it's already saved
	isEmailAlreadySaved := data.IsEmailAddressSaved(email)
	if isEmailAlreadySaved {
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprintln(w, "This email address is already saved.")
		return
	}

	// ensure it gets saved
	err := data.AddEmailAddress(email)
	if err != nil {
		log.Fatal("Email Not Added")
	}

	_, _ = fmt.Fprintln(w, "Success")
}

// sendEmailsHandler is responsible for POST requests at /sendEmails
func sendEmailsHandler(w http.ResponseWriter, _ *http.Request) {
	err := emails.SendRate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = fmt.Fprintln(w, "Success")
}
