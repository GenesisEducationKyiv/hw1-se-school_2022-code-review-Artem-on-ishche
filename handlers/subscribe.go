package handlers

import (
	"errors"
	"net/http"

	"gses2.app/api/services"
)

var EmailAddressesStorageImpl services.EmailAddressesStorage
var SubscribeRequestHandler = subscribeRequestHandler{}

type subscribeRequestHandler struct{}

var errMissingParameter = errors.New("required parameter is missing")

func (h subscribeRequestHandler) HandleRequest(request *http.Request) httpResponse {
	emailAddressString, err := getEmailParameter(request)
	if err != nil {
		return newHttpResponse(http.StatusBadRequest, "Required parameter 'email' is missing")
	}

	emailAddress, err := services.NewEmailAddress(emailAddressString)
	if err != nil {
		return newHttpResponse(http.StatusBadRequest, "Provided email address is wrong")
	}

	err = services.AddEmailAddress(EmailAddressesStorageImpl, *emailAddress)
	if isEmailAlreadySaved(err, emailAddressString) {
		return newHttpResponse(http.StatusConflict, "This email address is already saved")
	} else if err != nil {
		return newHttpResponse(http.StatusBadRequest, "Error when saving the email address")
	}

	return newHttpResponse(http.StatusOK, "Success")
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

func isEmailAlreadySaved(err error, emailAddressString string) bool {
	return err == services.ErrEmailAddressAlreadyExists(emailAddressString)
}
