package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	handleRequests()
}

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

func rateHandler(w http.ResponseWriter, _ *http.Request) {
	rate, err := getBtcUahRate()
	if err != nil || rate == -1 {
		http.Error(w, "Error has occurred", http.StatusBadRequest)
		return
	}
	_, _ = fmt.Fprint(w, rate)
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	// get the parameter from request
	emails, ok := r.URL.Query()["email"]

	// if the request doesn't have a required parameter
	if !ok || len(emails[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "Bad Request")
		return
	}

	// save only the first email if more are provided
	email := emails[0]

	// check if it's already saved
	isEmailAlreadySaved := isEmailAddressSaved(email)
	if isEmailAlreadySaved {
		w.WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprintln(w, "This email address is already saved.")
		return
	}

	// ensure it gets saved
	err := addEmailAddress(email)
	if err != nil {
		log.Fatal("Email Not Added")
	}

	_, _ = fmt.Fprintln(w, "Success")
}

func sendEmailsHandler(w http.ResponseWriter, _ *http.Request) {
	sendRate()
	_, _ = fmt.Fprintln(w, "Success")
}
