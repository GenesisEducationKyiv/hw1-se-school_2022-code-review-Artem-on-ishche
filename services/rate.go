package services

import "errors"

var ErrAPIRequestUnsuccessful = errors.New("API request has been unsuccessful")
var ErrAPIResponseUnmarshallError = errors.New("error when unmarshalling API response")

type ExchangeRateService interface {
	GetExchangeRate(pair CurrencyPair) (float64, error)
}

type BtcToUahRateService interface {
	GetBtcToUahRate() (float64, error)
}

type btcToUahRateServiceImpl struct {
	genericRateService ExchangeRateService
}

func NewBtcToUahServiceImpl(genericRateService ExchangeRateService) BtcToUahRateService {
	return &btcToUahRateServiceImpl{genericRateService}
}

func (btcUahService *btcToUahRateServiceImpl) GetBtcToUahRate() (float64, error) {
	btcUahPair := NewCurrencyPair(NewCurrency("BTC"), NewCurrency("UAH"))

	return btcUahService.genericRateService.GetExchangeRate(btcUahPair)
}
