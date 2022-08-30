package handlers

import (
	"errors"
	"log"
	"net/http"

	"gses2.app/api/data"
)

var errMissingParameter = errors.New("required parameter is missing")

func SubscribeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	email, err := getEmailParameter(request)
	if err != nil {
		sendBadRequestResponse(responseWriter, "Bad Request")

		return
	}

	isEmailAlreadySaved := data.IsEmailAddressSaved(email)
	if isEmailAlreadySaved {
		sendConflictResponse(responseWriter, "This email address is already saved")

		return
	}

	err = data.AddEmailAddress(email)
	if err != nil {
		log.Fatal("Failed to add an email")
	}

	sendSuccessResponse(responseWriter, "Success")
}

func getEmailParameter(request *http.Request) (string, error) {
	emailParams, ok := request.URL.Query()["email"]
	if !hasRequiredParameter(emailParams, ok) {
		return "", errMissingParameter
	}

	// save only the first email if more are provided
	return emailParams[0], nil
}

func hasRequiredParameter(emailParams []string, ok bool) bool {
	return ok && len(emailParams[0]) > 0
}
