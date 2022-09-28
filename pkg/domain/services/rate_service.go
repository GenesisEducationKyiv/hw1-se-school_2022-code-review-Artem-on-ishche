package services

import (
	"errors"

	"gses2.app/api/pkg/domain/models"
)

var (
	ErrAPIRequestUnsuccessful     = errors.New("API request has been unsuccessful")
	ErrAPIResponseUnmarshallError = errors.New("error when unmarshalling API response")
)

type ExchangeRateService interface {
	GetExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error)
}
