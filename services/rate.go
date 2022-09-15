package services

import "errors"

var ErrAPIRequestUnsuccessful = errors.New("API request has been unsuccessful")
var ErrAPIResponseUnmarshallError = errors.New("error when unmarshalling API response")

type ExchangeRateService interface {
	GetExchangeRate(from, to Currency) (float64, error)
	SetNext(service *ExchangeRateService)
}

type RateServiceFactory interface {
	CreateRateService() ExchangeRateService
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
