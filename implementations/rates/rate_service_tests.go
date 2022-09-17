package rates

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/services"
)

var (
	oneBillion        = 1_000_000_000
	oneSecondDuration = time.Duration(oneBillion)
)

var tests = []func(t *testing.T, service services.ExchangeRateService){
	testThatGetExchangeRateReturnsErrorForUnsupportedCurrencies,
	testThatGetExchangeRateReturnsMinusOneForUnsupportedCurrencies,
	testThatGetExchangeRateReturnsNoErrorForSupportedCurrencies,
	testThatGetExchangeRateReturnValuesDontFluctuateMuchOnSuccessiveCallsAfterOneSecond,
}

func testRateAPIClient(t *testing.T, client services.ExchangeRateService) {
	t.Helper()

	for _, test := range tests {
		test(t, client)
		time.Sleep(oneSecondDuration)
	}
}

func testThatGetExchangeRateReturnsErrorForUnsupportedCurrencies(t *testing.T, service services.ExchangeRateService) {
	t.Helper()

	currencyPair := getUnsupportedCurrencies()

	_, err := service.GetExchangeRate(currencyPair)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, services.ErrAPIRequestUnsuccessful)
}

func testThatGetExchangeRateReturnsMinusOneForUnsupportedCurrencies(t *testing.T, service services.ExchangeRateService) {
	t.Helper()

	currencyPair := getUnsupportedCurrencies()

	rate, _ := service.GetExchangeRate(currencyPair)

	assert.Equal(t, float64(-1), rate)
}

func testThatGetExchangeRateReturnsNoErrorForSupportedCurrencies(t *testing.T, service services.ExchangeRateService) {
	t.Helper()

	currencyPair := getSupportedCurrencies()

	_, err := service.GetExchangeRate(currencyPair)

	assert.Nil(t, err)
}

func testThatGetExchangeRateReturnValuesDontFluctuateMuchOnSuccessiveCallsAfterOneSecond(
	t *testing.T, service services.ExchangeRateService,
) {
	t.Helper()

	currencyPair := getSupportedCurrencies()
	deltaFraction := 0.25

	rate1, err1 := service.GetExchangeRate(currencyPair)

	time.Sleep(oneSecondDuration)

	rate2, err2 := service.GetExchangeRate(currencyPair)

	delta := rate1 * deltaFraction

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.InDelta(t, rate1, rate2, delta)
}

func getUnsupportedCurrencies() services.CurrencyPair {
	return services.NewCurrencyPair(services.NewCurrency("aaaaaaaaaa"), services.NewCurrency("bbbbbbbbbb"))
}

func getSupportedCurrencies() services.CurrencyPair {
	return services.NewCurrencyPair(services.NewCurrency("btc"), services.NewCurrency("uah"))
}
