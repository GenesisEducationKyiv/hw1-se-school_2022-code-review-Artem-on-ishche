package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func sendBadRequestResponse(responseWriter http.ResponseWriter, message string) {
	http.Error(responseWriter, message, http.StatusBadRequest)
}

func sendConflictResponse(responseWriter http.ResponseWriter, message string) {
	http.Error(responseWriter, message, http.StatusConflict)
}

func sendInternalServerErrorResponse(responseWriter http.ResponseWriter) {
	responseWriter.WriteHeader(http.StatusInternalServerError)
}

func sendSuccessResponse(responseWriter http.ResponseWriter, message string) {
	_, err := fmt.Fprintln(responseWriter, message)
	if err != nil {
		log.Println("Error when sending a success response")
	}
}
