package http

import (
	"errors"
	"net/http"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

var errMissingParameter = errors.New("required parameter is missing")

type SubscribeRequestHandler struct {
	AddEmailAddressService application.AddEmailAddressService
}

func (handler SubscribeRequestHandler) GetPath() string {
	return "/subscribe"
}

func (handler SubscribeRequestHandler) GetMethod() string {
	return "POST"
}

func (handler SubscribeRequestHandler) GetResponse(request *http.Request) Response {
	emailAddressString, err := getEmailParameter(request)
	if err != nil {
		return newResponse(http.StatusBadRequest, "Required parameter 'email' is missing")
	}

	emailAddress, err := models.NewEmailAddress(emailAddressString)
	if err != nil {
		return newResponse(http.StatusBadRequest, "Provided email address is wrong")
	}

	err = handler.AddEmailAddressService.AddEmailAddress(*emailAddress)
	if err == nil {
		return newResponse(http.StatusOK, "Success")
	} else if isEmailAlreadySaved(err, emailAddressString) {
		return newResponse(http.StatusConflict, "This email address is already saved")
	} else {
		return newResponse(http.StatusInternalServerError, "Error when saving the email address")
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
