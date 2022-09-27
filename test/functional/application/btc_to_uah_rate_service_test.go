package application

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type (
	getRateFunction               func(pair models.CurrencyPair) (float64, error)
	getRateFunctionReturnedValues struct {
		rate float64
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

func (rateService *exchangeRateServiceTestDouble) GetExchangeRate(pair models.CurrencyPair) (float64, error) {
	return getRateTestFunction(pair)
}

var generalGetRateTests = []getRateTest{
	{
		function: func(pair models.CurrencyPair) (float64, error) {
			return 100, nil
		},
		expectedReturn: getRateFunctionReturnedValues{100, nil},
	},
	{
		function: func(pair models.CurrencyPair) (float64, error) {
			return -1, services.ErrAPIRequestUnsuccessful
		},
		expectedReturn: getRateFunctionReturnedValues{-1, services.ErrAPIRequestUnsuccessful},
	},
	{
		function: func(pair models.CurrencyPair) (float64, error) {
			if pair.From.Name == "BTC" && pair.To.Name == "UAH" {
				return 900_000.001, nil
			}

			return -1, nil
		},
		expectedReturn: getRateFunctionReturnedValues{900000.001, nil},
	},
	{
		function: func(pair models.CurrencyPair) (float64, error) {
			if pair.From.Name == "BTC" && pair.To.Name == "UAH" {
				return -1, services.ErrAPIRequestUnsuccessful
			}

			return 3742.134, nil
		},
		expectedReturn: getRateFunctionReturnedValues{-1, services.ErrAPIRequestUnsuccessful},
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

	getRateTestFunction = func(pair models.CurrencyPair) (float64, error) {
		callsCount++

		fromCurrencyName = pair.From.Name
		toCurrencyName = pair.To.Name

		return 1, nil
	}

	_, _ = btcToUahServiceImpl.GetBtcToUahRate()

	assert.Equal(t, 1, callsCount)
	assert.Equal(t, "BTC", fromCurrencyName)
	assert.Equal(t, "UAH", toCurrencyName)
}
