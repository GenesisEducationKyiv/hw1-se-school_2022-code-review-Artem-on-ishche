package rates

import (
	"gopkg.in/resty.v0"
	"gses2.app/api/services"
	"time"
)

type parsedResponse struct {
	rate float64
	time time.Time
}

type exchangeRateAPIClient interface {
	getName() string
	getAPIRequestUrlForGivenCurrencies(pair services.CurrencyPair) string
	getAPIRequest() *resty.Request
	parseResponseBody(responseBody []byte) (*parsedResponse, error)
}
