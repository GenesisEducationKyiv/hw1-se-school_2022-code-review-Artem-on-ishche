package application

import "gses2.app/api/pkg/domain/models"

var btcUahPair = models.CurrencyPair{
	Base:  models.Currency{Name: "BTC"},
	Quote: models.Currency{Name: "UAH"},
}

type getRateFunction func(pair models.CurrencyPair) (*models.ExchangeRate, error)

var getRateTestFunction getRateFunction

type rateServiceTestDouble struct{}

func (s rateServiceTestDouble) GetExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error) {
	return getRateTestFunction(pair)
}
