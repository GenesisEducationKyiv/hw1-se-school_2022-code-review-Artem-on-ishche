package rates

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/resty.v0"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
)

const binanceRequestFormatString = "https://api.binance.com/api/v3/avgPrice?symbol=%s%s"

type receivedBinanceAPIResponse struct {
	Mins  string `json:"mins"`
	Price string `json:"price"`
}

type BinanceAPIClientFactory struct {
	Mediator *Mediator
}

func (factory BinanceAPIClientFactory) CreateRateService() ExchangeRateServiceChain {
	return &exchangeRateService{
		mediator:           factory.Mediator,
		concreteRateClient: &binanceAPIClient{},
	}
}

type binanceAPIClient struct{}

func (c *binanceAPIClient) name() string {
	return "Binance"
}

func (c *binanceAPIClient) getAPIRequestURLForGivenCurrencies(pair models.CurrencyPair) string {
	return fmt.Sprintf(
		binanceRequestFormatString,
		pair.Base.Name,
		pair.Quote.Name,
	)
}

func (c *binanceAPIClient) getAPIRequest() *resty.Request {
	return resty.R()
}

func (c *binanceAPIClient) parseResponseBody(responseBody []byte) (*parsedResponse, error) {
	var result receivedBinanceAPIResponse

	err := json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, application.ErrAPIResponseUnmarshallError
	}

	bitSize := 64

	price, err := strconv.ParseFloat(result.Price, bitSize)
	if err != nil {
		return nil, application.ErrAPIRequestUnsuccessful
	}

	return &parsedResponse{
		price: price,
		time:  time.Now(),
	}, nil
}
