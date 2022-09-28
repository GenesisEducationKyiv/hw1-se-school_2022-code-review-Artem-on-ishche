package rates

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

var (
	oneBillion        = 1_000_000_000
	oneSecondDuration = time.Duration(oneBillion)
)

var tests = []func(t *testing.T, service services.ExchangeRateService){
	testThatGetExchangeRateReturnsErrorForUnsupportedCurrencies,
	testThatGetExchangeRateReturnsNilForUnsupportedCurrencies,
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

func testThatGetExchangeRateReturnsNilForUnsupportedCurrencies(t *testing.T, service services.ExchangeRateService) {
	t.Helper()

	currencyPair := getUnsupportedCurrencies()

	rate, _ := service.GetExchangeRate(currencyPair)

	assert.Equal(t, nil, rate)
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

	delta := rate1.Price * deltaFraction

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.InDelta(t, rate1.Price, rate2.Price, delta)
}

func getUnsupportedCurrencies() models.CurrencyPair {
	return models.NewCurrencyPair(models.NewCurrency("aaaaaaaaaa"), models.NewCurrency("bbbbbbbbbb"))
}

func getSupportedCurrencies() models.CurrencyPair {
	return models.NewCurrencyPair(models.NewCurrency("btc"), models.NewCurrency("uah"))
}
