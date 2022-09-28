package application

import (
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type BtcToUahRateService interface {
	GetBtcToUahRate() (*models.ExchangeRate, error)
}

type btcToUahRateServiceImpl struct {
	genericRateService services.ExchangeRateService
}

func NewBtcToUahServiceImpl(genericRateService services.ExchangeRateService) BtcToUahRateService {
	return &btcToUahRateServiceImpl{genericRateService}
}

func (btcUahService *btcToUahRateServiceImpl) GetBtcToUahRate() (*models.ExchangeRate, error) {
	btcUahPair := models.NewCurrencyPair(models.NewCurrency("BTC"), models.NewCurrency("UAH"))

	return btcUahService.genericRateService.GetExchangeRate(btcUahPair)
}
