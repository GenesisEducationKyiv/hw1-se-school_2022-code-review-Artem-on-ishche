package http

import (
	"errors"
	"fmt"
	"net/http"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/services"
)

type BtcToUahRateRequestHandler struct {
	BtcToUahService application.BtcToUahRateService
}

func (handler BtcToUahRateRequestHandler) GetPath() string {
	return "/rate"
}

func (handler BtcToUahRateRequestHandler) GetMethod() string {
	return "GET"
}

func (handler BtcToUahRateRequestHandler) GetResponse(_ *http.Request) Response {
	exchangeRate, err := handler.BtcToUahService.GetBtcToUahRate()
	if errors.Is(err, nil) {
		return newResponse(http.StatusOK, fmt.Sprintf("%v", exchangeRate.Price))
	} else if errors.Is(err, services.ErrAPIRequestUnsuccessful) {
		return newResponse(http.StatusBadGateway, "API request has not been successful")
	} else {
		return newResponse(http.StatusInternalServerError, "Some unexpected error has occurred")
	}
}
