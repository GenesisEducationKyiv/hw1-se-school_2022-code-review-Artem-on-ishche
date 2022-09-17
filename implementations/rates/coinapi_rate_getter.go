package rates

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v0"
	"gses2.app/api/config"
	"gses2.app/api/services"
	"time"
)

type receivedCoinAPIResponse struct {
	Time string  `json:"time"`
	Rate float64 `json:"rate"`
}

type CoinAPIClientFactory struct {
	Mediator *Mediator
}

func (factory CoinAPIClientFactory) CreateRateService() ExchangeRateServiceChain {
	return &exchangeRateService{
		mediator:           factory.Mediator,
		concreteRateClient: coinAPIClient{},
	}
}

type coinAPIClient struct{}

func (c coinAPIClient) getName() string {
	return "Coinbase"
}

func (c coinAPIClient) getAPIRequestUrlForGivenCurrencies(pair services.CurrencyPair) string {
	return fmt.Sprintf("https://rest.coinapi.io/v1/exchangerate/%v/%v", pair.From.Name, pair.To.Name)
}

func (c coinAPIClient) getAPIRequest() *resty.Request {
	return resty.R().SetHeader("X-CoinAPI-Key", config.CoinAPIKeyValue)
}

func (c coinAPIClient) parseResponseBody(responseBody []byte) (*parsedResponse, error) {
	var result receivedCoinAPIResponse

	err := json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, services.ErrAPIResponseUnmarshallError
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05.999Z", result.Time)
	if err != nil {
		return nil, services.ErrAPIResponseUnmarshallError
	}

	return &parsedResponse{
		rate: result.Rate,
		time: timestamp,
	}, nil
}
