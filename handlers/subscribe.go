package handlers

import (
	"errors"
	"net/http"

	"gses2.app/api/services"
)

var errMissingParameter = errors.New("required parameter is missing")

type subscribeRequestHandler struct {
	addEmailAddressService services.AddEmailAddressService
}

func NewSubscribeRequestHandler(addEmailAddressService services.AddEmailAddressService) RequestHandler {
	return subscribeRequestHandler{addEmailAddressService}
}

func (handler subscribeRequestHandler) HandleRequest(request *http.Request) httpResponse {
	emailAddressString, err := getEmailParameter(request)
	if err != nil {
		return newHTTPResponse(http.StatusBadRequest, "Required parameter 'email' is missing")
	}

	emailAddress, err := services.NewEmailAddress(emailAddressString)
	if err != nil {
		return newHTTPResponse(http.StatusBadRequest, "Provided email address is wrong")
	}

	err = handler.addEmailAddressService.AddEmailAddress(*emailAddress)
	if err == nil {
		return newHTTPResponse(http.StatusOK, "Success")
	} else if isEmailAlreadySaved(err, emailAddressString) {
		return newHTTPResponse(http.StatusConflict, "This email address is already saved")
	} else {
		return newHTTPResponse(http.StatusInternalServerError, "Error when saving the email address")
	}
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
	return err.Error() == services.ErrEmailAddressAlreadyExists(emailAddressString).Error()
}
