package handlers

import (
	"fmt"
	"net/http"

	"gses2.app/api/services"
)

var ExchangeRateServiceImpl services.ExchangeRateService
var RateRequestHandler = rateRequestHandler{}

type rateRequestHandler struct{}

func (h rateRequestHandler) HandleRequest(_ *http.Request) httpResponse {
	exchangeRate, err := services.GetBtcToUahRate(ExchangeRateServiceImpl)

	if !isApiRequestSuccessful(err) {
		return newHttpResponse(http.StatusBadRequest, "API request has not been successful")
	} else if isRateWrong(exchangeRate, err) {
		return newHttpResponse(http.StatusBadRequest, "Some unexpected error has occurred")
	} else {
		rateString := fmt.Sprintf("%v", exchangeRate)
		return newHttpResponse(http.StatusOK, rateString)
	}
}

func isApiRequestSuccessful(err error) bool {
	return err != services.ErrApiRequestUnsuccessful
}

func isRateWrong(rate float64, err error) bool {
	return err != nil || rate <= 0
}
