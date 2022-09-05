package implementations

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/config"
	"gses2.app/api/services"
)

func TestThatGetExchangeRateReturnsErrorForUnsupportedCurrencies(t *testing.T) {
	config.LoadEnv()
	currencyFrom, currencyTo := getUnsupportedCurrencies()

	_, err := GetExchangeRateCoinApiClient().GetExchangeRate(currencyFrom, currencyTo)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, services.ErrApiRequestUnsuccessful)
}

func TestThatGetExchangeRateReturnsMinusOneForUnsupportedCurrencies(t *testing.T) {
	config.LoadEnv()
	currencyFrom, currencyTo := getUnsupportedCurrencies()

	rate, _ := GetExchangeRateCoinApiClient().GetExchangeRate(currencyFrom, currencyTo)

	assert.Equal(t, float64(-1), rate)
}

func TestThatGetExchangeRateReturnsNoErrorForSupportedCurrencies(t *testing.T) {
	config.LoadEnv()
	currencyFrom, currencyTo := getSupportedCurrencies()

	_, err := GetExchangeRateCoinApiClient().GetExchangeRate(currencyFrom, currencyTo)

	assert.Nil(t, err)
}

func TestThatGetExchangeRateReturnValuesDontFluctuateMuchOnSuccessiveCallsAfterOneSecond(t *testing.T) {
	config.LoadEnv()
	currencyFrom, currencyTo := getSupportedCurrencies()
	oneSecondDuration := time.Duration(1_000_000_000)

	rate1, err1 := GetExchangeRateCoinApiClient().GetExchangeRate(currencyFrom, currencyTo)
	time.Sleep(oneSecondDuration)
	rate2, err2 := GetExchangeRateCoinApiClient().GetExchangeRate(currencyFrom, currencyTo)
	delta := rate1 * 0.25 // delta is 25% of rate1

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.InDelta(t, rate1, rate2, delta)
}

func getUnsupportedCurrencies() (services.Currency, services.Currency) {
	return services.NewCurrency("AAA"), services.NewCurrency("BBB")
}

func getSupportedCurrencies() (services.Currency, services.Currency) {
	return services.NewCurrency("uah"), services.NewCurrency("btc")
}
