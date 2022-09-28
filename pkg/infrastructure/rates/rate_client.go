package rates

import (
	"time"

	"gopkg.in/resty.v0"

	"gses2.app/api/pkg/domain/models"
)

type parsedResponse struct {
	price float64
	time  time.Time
}

type exchangeRateAPIClient interface {
	getName() string
	getAPIRequestURLForGivenCurrencies(pair models.CurrencyPair) string
	getAPIRequest() *resty.Request
	parseResponseBody(responseBody []byte) (*parsedResponse, error)
}
