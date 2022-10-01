package http

import (
	"gses2.app/api/pkg/domain/models"
	httpPresentation "gses2.app/api/pkg/presentation/http/routes"
)

type getRateFunction func() (*models.ExchangeRate, error)

var getRateTestFunction getRateFunction

type rateServiceTestDouble struct{}

func (service rateServiceTestDouble) GetExchangeRate(models.CurrencyPair) (*models.ExchangeRate, error) {
	return getRateTestFunction()
}

var testBtcToUahHandler = httpPresentation.RateRoute{ExchangeRateService: rateServiceTestDouble{}}
