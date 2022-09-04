package handlers

import (
	"net/http"

	"gses2.app/api/services"
)

var EmailSenderImpl services.EmailSender
var SendEmailsRequestHandler = sendEmailsRequestHandler{}

type sendEmailsRequestHandler struct{}

func (h sendEmailsRequestHandler) HandleRequest(request *http.Request) httpResponse {
	err := services.SendBtcToUahRateEmails(ExchangeRateServiceImpl, EmailAddressesStorageImpl, EmailSenderImpl)

	switch err {
	case nil:
		return newHttpResponse(http.StatusOK, "Success")
	case services.ErrApiRequestUnsuccessful:
		return newHttpResponse(http.StatusBadRequest, "API request has not been successful")
	default:
		return newHttpResponse(http.StatusBadRequest, "Some error has occurred")
	}
}
