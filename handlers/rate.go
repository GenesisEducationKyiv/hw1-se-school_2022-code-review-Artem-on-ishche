package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"gses2.app/api/services"
)

type btcToUahRateRequestHandler struct {
	btcToUahService services.BtcToUahRateService
}

func NewBtcToUahRateRequestHandler(genericExchangeRateService services.ExchangeRateService) RequestHandler {
	btcToUahRateService := services.NewBtcToUahServiceImpl(genericExchangeRateService)

	return btcToUahRateRequestHandler{btcToUahRateService}
}

func (handler btcToUahRateRequestHandler) HandleRequest(_ *http.Request) httpResponse {
	exchangeRate, err := handler.btcToUahService.GetBtcToUahRate()
	if errors.Is(err, nil) {
		return newHTTPResponse(http.StatusOK, fmt.Sprintf("%v", exchangeRate))
	} else if errors.Is(err, services.ErrAPIRequestUnsuccessful) {
		return newHTTPResponse(http.StatusBadGateway, "API request has not been successful")
	} else {
		return newHTTPResponse(http.StatusInternalServerError, "Some unexpected error has occurred")
	}
}
