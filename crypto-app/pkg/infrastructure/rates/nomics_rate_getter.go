package rates

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/resty.v0"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

const nomicsRequestFormatString = "https://api.nomics.com/v1/currencies/ticker?key=%v&ids=%v&interval=1d&convert=%v"

type receivedNomicsAPIResponse struct {
	Price          string `json:"price"`
	PriceTimestamp string `json:"price_timestamp"`
}

type NomicsAPIClientFactory struct {
	Cacher CacherRateService
	Logger services.Logger
}

func (factory NomicsAPIClientFactory) CreateRateService() ExchangeRateServiceChain {
	return &exchangeRateService{
		concreteRateClient: nomicsAPIClient{},
		cacher:             factory.Cacher,
		logger:             factory.Logger,
	}
}

type nomicsAPIClient struct{}

func (c nomicsAPIClient) name() string {
	return "Nomics"
}

func (c nomicsAPIClient) getAPIRequestURLForGivenCurrencies(pair models.CurrencyPair) string {
	return fmt.Sprintf(
		nomicsRequestFormatString,
		config.NomicsAPIKeyValue,
		pair.Base.Name,
		pair.Quote.Name,
	)
}

func (c nomicsAPIClient) getAPIRequest() *resty.Request {
	return resty.R()
}

func (c nomicsAPIClient) parseResponseBody(responseBody []byte) (*parsedResponse, error) {
	var results []receivedNomicsAPIResponse

	err := json.Unmarshal(responseBody, &results)
	if err != nil {
		return nil, application.ErrAPIResponseUnmarshallError
	}

	if len(results) == 0 {
		return nil, application.ErrAPIRequestUnsuccessful
	}

	result := results[0]
	bitSize := 64

	price, err := strconv.ParseFloat(result.Price, bitSize)
	if err != nil {
		return nil, application.ErrAPIRequestUnsuccessful
	}

	timestamp, err := time.Parse(timeLayout, result.PriceTimestamp)
	if err != nil {
		return nil, application.ErrAPIResponseUnmarshallError
	}

	return &parsedResponse{
		price: price,
		time:  timestamp,
	}, nil
}
