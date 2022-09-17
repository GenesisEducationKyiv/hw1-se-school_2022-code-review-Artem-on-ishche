package rates

import (
	"time"

	"gopkg.in/resty.v0"

	"gses2.app/api/services"
)

type parsedResponse struct {
	rate float64
	time time.Time
}

type exchangeRateAPIClient interface {
	getName() string
	getAPIRequestURLForGivenCurrencies(pair services.CurrencyPair) string
	getAPIRequest() *resty.Request
	parseResponseBody(responseBody []byte) (*parsedResponse, error)
}
