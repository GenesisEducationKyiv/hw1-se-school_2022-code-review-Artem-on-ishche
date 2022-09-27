package http

import (
	"errors"
	"net/http"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/services"
)

type SendEmailsRequestHandler struct {
	SendBtcToUahRateEmailsService application.SendBtcToUahRateEmailsService
}

func (handler SendEmailsRequestHandler) GetPath() string {
	return "/sendEmails"
}

func (handler SendEmailsRequestHandler) GetMethod() string {
	return "POST"
}

func (handler SendEmailsRequestHandler) GetResponse(request *http.Request) Response {
	err := handler.SendBtcToUahRateEmailsService.SendBtcToUahRateEmails()
	if errors.Is(err, nil) {
		return newResponse(http.StatusOK, "Success")
	} else if errors.Is(err, services.ErrAPIRequestUnsuccessful) {
		return newResponse(http.StatusBadGateway, "API request has not been successful")
	} else {
		return newResponse(http.StatusInternalServerError, "Some error has occurred")
	}
}
