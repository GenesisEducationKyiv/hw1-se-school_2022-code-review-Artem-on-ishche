package services

import "errors"

var ErrApiRequestUnsuccessful = errors.New("API request has been unsuccessful")

type ExchangeRateService interface {
	GetExchangeRate(from, to Currency) (float64, error)
}

func GetBtcToUahRate(service ExchangeRateService) (float64, error) {
	return service.GetExchangeRate(NewCurrency("BTC"), NewCurrency("UAH"))
}
