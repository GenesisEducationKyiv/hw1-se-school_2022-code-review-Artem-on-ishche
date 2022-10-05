package services

import (
	"gses2.app/api/pkg/domain/models"
)

type ExchangeRateService interface {
	GetExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error)
}
