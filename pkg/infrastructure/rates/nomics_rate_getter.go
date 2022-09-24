package rates

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/resty.v0"

	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type receivedNomicsAPIResponse struct {
	Price          string `json:"price"`
	PriceTimestamp string `json:"price_timestamp"`
}

type NomicsAPIClientFactory struct {
	Mediator *Mediator
}

func (factory NomicsAPIClientFactory) CreateRateService() ExchangeRateServiceChain {
	return &exchangeRateService{
		mediator:           factory.Mediator,
		concreteRateClient: nomicsAPIClient{},
	}
}

type nomicsAPIClient struct{}

func (c nomicsAPIClient) getName() string {
	return "Nomics"
}

func (c nomicsAPIClient) getAPIRequestURLForGivenCurrencies(pair models.CurrencyPair) string {
	return fmt.Sprintf(
		"https://api.nomics.com/v1/currencies/ticker?key=%v&ids=%v&interval=1d&convert=%v",
		config.NomicsAPIKeyValue,
		pair.From.Name,
		pair.To.Name,
	)
}

func (c nomicsAPIClient) getAPIRequest() *resty.Request {
	return resty.R()
}

func (c nomicsAPIClient) parseResponseBody(responseBody []byte) (*parsedResponse, error) {
	var results []receivedNomicsAPIResponse

	err := json.Unmarshal(responseBody, &results)
	if err != nil {
		return nil, services.ErrAPIResponseUnmarshallError
	}

	if len(results) == 0 {
		return nil, services.ErrAPIRequestUnsuccessful
	}

	result := results[0]
	bitSize := 64

	price, err := strconv.ParseFloat(result.Price, bitSize)
	if err != nil {
		return nil, services.ErrAPIRequestUnsuccessful
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05Z", result.PriceTimestamp)
	if err != nil {
		return nil, services.ErrAPIResponseUnmarshallError
	}

	return &parsedResponse{
		rate: price,
		time: timestamp,
	}, nil
}
