package handlers

import (
	"errors"
	"net/http"

	"gses2.app/api/services"
)

type sendEmailsRequestHandler struct {
	sendBtcToUahRateEmailsService services.SendBtcToUahRateEmailsService
}

func NewSendEmailsRequestHandler(
	genericExchangeRateService services.ExchangeRateService,
	emailAddressesStorage services.EmailAddressesStorage,
	emailSender services.EmailSender,
) RequestHandler {
	btcToUahRateService := services.NewBtcToUahServiceImpl(genericExchangeRateService)
	sendBtcToUahRateEmailsService :=
		services.NewSendBtcToUahRateEmailsServiceImpl(btcToUahRateService, emailAddressesStorage, emailSender)

	return sendEmailsRequestHandler{sendBtcToUahRateEmailsService}
}

func (handler sendEmailsRequestHandler) HandleRequest(request *http.Request) httpResponse {
	err := handler.sendBtcToUahRateEmailsService.SendBtcToUahRateEmails()
	if errors.Is(err, nil) {
		return newHTTPResponse(http.StatusOK, "Success")
	} else if errors.Is(err, services.ErrAPIRequestUnsuccessful) {
		return newHTTPResponse(http.StatusBadGateway, "API request has not been successful")
	} else {
		return newHTTPResponse(http.StatusInternalServerError, "Some error has occurred")
	}
}
