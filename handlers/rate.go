package handlers

import (
	"fmt"
	"net/http"

	"gses2.app/api/services"
)

type btcToUahRateRequestHandler struct {
	btcToUahService services.BtcToUahRateService
}

func NewBtcToUahRateRequestHandler(btcToUahService services.BtcToUahRateService) RequestHandler {
	return btcToUahRateRequestHandler{btcToUahService}
}

func (handler btcToUahRateRequestHandler) HandleRequest(_ *http.Request) httpResponse {
	exchangeRate, err := handler.btcToUahService.GetBtcToUahRate()

	switch err {
	case nil:
		return newHttpResponse(http.StatusOK, fmt.Sprintf("%v", exchangeRate))
	case services.ErrApiRequestUnsuccessful:
		return newHttpResponse(http.StatusBadGateway, "API request has not been successful")
	default:
		return newHttpResponse(http.StatusInternalServerError, "Some unexpected error has occurred")
	}
}
