package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	getRateFunction               func(from, to Currency) (float64, error)
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

func (rateService *exchangeRateServiceTestDouble) GetExchangeRate(from, to Currency) (float64, error) {
	return getRateTestFunction(from, to)
}

var generalGetRateTests = []getRateTest{
	{
		function: func(from, to Currency) (float64, error) {
			return 100, nil
		},
		expectedReturn: getRateFunctionReturnedValues{100, nil},
	},
	{
		function: func(from, to Currency) (float64, error) {
			return -1, ErrAPIRequestUnsuccessful
		},
		expectedReturn: getRateFunctionReturnedValues{-1, ErrAPIRequestUnsuccessful},
	},
	{
		function: func(from, to Currency) (float64, error) {
			if from.Name == "BTC" && to.Name == "UAH" {
				return 900_000.001, nil
			}

			return -1, nil
		},
		expectedReturn: getRateFunctionReturnedValues{900000.001, nil},
	},
	{
		function: func(from, to Currency) (float64, error) {
			if from.Name == "BTC" && to.Name == "UAH" {
				return -1, ErrAPIRequestUnsuccessful
			}

			return 3742.134, nil
		},
		expectedReturn: getRateFunctionReturnedValues{-1, ErrAPIRequestUnsuccessful},
	},
}

func TestThatGetBtcToUahRateReturnsCorrectValues(t *testing.T) {
	fakeService := exchangeRateServiceTestDouble{}
	btcToUahServiceImpl := NewBtcToUahServiceImpl(&fakeService)

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
	btcToUahServiceImpl := NewBtcToUahServiceImpl(&spyService)

	callsCount := 0

	var fromCurrencyName, toCurrencyName string

	getRateTestFunction = func(from, to Currency) (float64, error) {
		callsCount++

		fromCurrencyName = from.Name
		toCurrencyName = to.Name

		return 1, nil
	}

	_, _ = btcToUahServiceImpl.GetBtcToUahRate()

	assert.Equal(t, 1, callsCount)
	assert.Equal(t, "BTC", fromCurrencyName)
	assert.Equal(t, "UAH", toCurrencyName)
}
