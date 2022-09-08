package services

import "errors"

var ErrAPIRequestUnsuccessful = errors.New("API request has been unsuccessful")

type ExchangeRateService interface {
	GetExchangeRate(from, to Currency) (float64, error)
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
	return btcUahService.genericRateService.GetExchangeRate(NewCurrency("BTC"), NewCurrency("UAH"))
}
