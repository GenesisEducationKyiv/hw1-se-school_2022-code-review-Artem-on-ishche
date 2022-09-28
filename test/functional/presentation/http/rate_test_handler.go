package http

import (
	"gses2.app/api/pkg/domain/models"
	httpPresentation "gses2.app/api/pkg/presentation/http"
)

type getRateFunction func() (*models.ExchangeRate, error)

var getRateTestFunction getRateFunction

type btcToUahServiceTestDouble struct{}

func (service btcToUahServiceTestDouble) GetBtcToUahRate() (*models.ExchangeRate, error) {
	return getRateTestFunction()
}

var testBtcToUahHandler = httpPresentation.BtcToUahRateRequestHandler{BtcToUahService: btcToUahServiceTestDouble{}}
