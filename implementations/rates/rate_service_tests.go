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

	currencyFrom, currencyTo := getUnsupportedCurrencies()

	_, err := service.GetExchangeRate(currencyFrom, currencyTo)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, services.ErrAPIRequestUnsuccessful)
}

func testThatGetExchangeRateReturnsMinusOneForUnsupportedCurrencies(t *testing.T, service services.ExchangeRateService) {
	t.Helper()

	currencyFrom, currencyTo := getUnsupportedCurrencies()

	rate, _ := service.GetExchangeRate(currencyFrom, currencyTo)

	assert.Equal(t, float64(-1), rate)
}

func testThatGetExchangeRateReturnsNoErrorForSupportedCurrencies(t *testing.T, service services.ExchangeRateService) {
	t.Helper()

	currencyFrom, currencyTo := getSupportedCurrencies()

	_, err := service.GetExchangeRate(currencyFrom, currencyTo)

	assert.Nil(t, err)
}

func testThatGetExchangeRateReturnValuesDontFluctuateMuchOnSuccessiveCallsAfterOneSecond(
	t *testing.T, service services.ExchangeRateService,
) {
	t.Helper()

	currencyFrom, currencyTo := getSupportedCurrencies()
	deltaFraction := 0.25

	rate1, err1 := service.GetExchangeRate(currencyFrom, currencyTo)

	time.Sleep(oneSecondDuration)

	rate2, err2 := service.GetExchangeRate(currencyFrom, currencyTo)

	delta := rate1 * deltaFraction

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.InDelta(t, rate1, rate2, delta)
}

func getUnsupportedCurrencies() (services.Currency, services.Currency) {
	return services.NewCurrency("aaaaaaaaaa"), services.NewCurrency("bbbbbbbbbb")
}

func getSupportedCurrencies() (services.Currency, services.Currency) {
	return services.NewCurrency("btc"), services.NewCurrency("uah")
}
