package application

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type (
	getRateFunction               func(pair models.CurrencyPair) (*models.ExchangeRate, error)
	getRateFunctionReturnedValues struct {
		rate *models.ExchangeRate
		err  error
	}
)

type getRateTest struct {
	function       getRateFunction
	expectedReturn getRateFunctionReturnedValues
}

var getRateTestFunction getRateFunction

type exchangeRateServiceTestDouble struct{}

func (rateService *exchangeRateServiceTestDouble) SetNext(_ *services.ExchangeRateService) {}

func (rateService *exchangeRateServiceTestDouble) GetExchangeRate(pair models.CurrencyPair) (*models.ExchangeRate, error) {
	return getRateTestFunction(pair)
}

var generalGetRateTests = []getRateTest{
	{
		function: func(pair models.CurrencyPair) (*models.ExchangeRate, error) {
			return models.NewExchangeRate(pair, 100), nil
		},
		expectedReturn: getRateFunctionReturnedValues{
			models.NewExchangeRate(
				models.NewCurrencyPair(models.NewCurrency("btc"), models.NewCurrency("uah")),
				100,
			),
			nil,
		},
	},
	{
		function: func(pair models.CurrencyPair) (*models.ExchangeRate, error) {
			return nil, services.ErrAPIRequestUnsuccessful
		},
		expectedReturn: getRateFunctionReturnedValues{nil, services.ErrAPIRequestUnsuccessful},
	},
	{
		function: func(pair models.CurrencyPair) (*models.ExchangeRate, error) {
			if pair.Base.Name == "BTC" && pair.Quote.Name == "UAH" {
				return models.NewExchangeRate(pair, 900_000.001), nil
			}

			return models.NewExchangeRate(pair, 1), nil
		},
		expectedReturn: getRateFunctionReturnedValues{
			models.NewExchangeRate(
				models.NewCurrencyPair(models.NewCurrency("btc"), models.NewCurrency("uah")),
				900_000.001,
			),
			nil,
		},
	},
	{
		function: func(pair models.CurrencyPair) (*models.ExchangeRate, error) {
			if pair.Base.Name == "BTC" && pair.Quote.Name == "UAH" {
				return nil, services.ErrAPIRequestUnsuccessful
			}

			return models.NewExchangeRate(pair, 3742.134), nil
		},
		expectedReturn: getRateFunctionReturnedValues{nil, services.ErrAPIRequestUnsuccessful},
	},
}

func TestThatGetBtcToUahRateReturnsCorrectValues(t *testing.T) {
	fakeService := exchangeRateServiceTestDouble{}
	btcToUahServiceImpl := application.NewBtcToUahServiceImpl(&fakeService)

	for _, test := range generalGetRateTests {
		getRateTestFunction = test.function
		expectedResult := test.expectedReturn

		rate, err := btcToUahServiceImpl.GetBtcToUahRate()

		assert.Equal(t, expectedResult.rate, rate)
		assert.Equal(t, expectedResult.err, err)
	}
}

func TestThatGetBtcToUahCallsRateServiceWithCorrectParameters(t *testing.T) {
	spyService := exchangeRateServiceTestDouble{}
	btcToUahServiceImpl := application.NewBtcToUahServiceImpl(&spyService)

	callsCount := 0

	var fromCurrencyName, toCurrencyName string

	getRateTestFunction = func(pair models.CurrencyPair) (*models.ExchangeRate, error) {
		callsCount++

		fromCurrencyName = pair.Base.Name
		toCurrencyName = pair.Quote.Name

		return models.NewExchangeRate(pair, 1), nil
	}

	_, _ = btcToUahServiceImpl.GetBtcToUahRate()

	assert.Equal(t, 1, callsCount)
	assert.Equal(t, "BTC", fromCurrencyName)
	assert.Equal(t, "UAH", toCurrencyName)
}
