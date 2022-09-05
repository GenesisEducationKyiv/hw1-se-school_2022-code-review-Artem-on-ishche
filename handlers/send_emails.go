package handlers

import (
	"net/http"

	"gses2.app/api/services"
)

type sendEmailsRequestHandler struct {
	sendBtcToUahRateEmailsService services.SendBtcToUahRateEmailsService
}

func NewSendEmailsRequestHandler(sendBtcToUahRateEmailsService services.SendBtcToUahRateEmailsService) RequestHandler {
	return sendEmailsRequestHandler{sendBtcToUahRateEmailsService}
}

func (handler sendEmailsRequestHandler) HandleRequest(request *http.Request) httpResponse {
	err := handler.sendBtcToUahRateEmailsService.SendBtcToUahRateEmails()

	switch err {
	case nil:
		return newHttpResponse(http.StatusOK, "Success")
	case services.ErrApiRequestUnsuccessful:
		return newHttpResponse(http.StatusBadGateway, "API request has not been successful")
	default:
		return newHttpResponse(http.StatusInternalServerError, "Some error has occurred")
	}
}
