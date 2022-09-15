package rates

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/resty.v0"

	"gses2.app/api/config"
	"gses2.app/api/services"
)

type receivedCoinMarketCapAPIUahResponse struct {
	Data []struct {
		Quote struct {
			Uah struct {
				Price float64 `json:"price"`
				Time  string  `json:"last_updated"`
			} `json:"UAH"`
		} `json:"quote"`
	} `json:"data"`
}

type CoinMarketCapAPIClientFactory struct{}

func (factory CoinMarketCapAPIClientFactory) CreateRateService() services.ExchangeRateService {
	return &exchangeRateService{concreteRateClient: coinMarketCapAPIClient{}}
}

type coinMarketCapAPIClient struct{}

func (c coinMarketCapAPIClient) getName() string {
	return "CoinMarketCap"
}

func (c coinMarketCapAPIClient) getAPIRequestUrlForGivenCurrencies(from, to services.Currency) string {
	return fmt.Sprintf(
		"https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?convert=%v&limit=1&start=1",
		to.Name,
	)
}

func (c coinMarketCapAPIClient) getAPIRequest() *resty.Request {
	return resty.R().
		SetHeader("X-CMC_PRO_API_KEY", config.CoinMarketCapAPIKeyValue)
}

func (c coinMarketCapAPIClient) parseResponseBody(responseBody []byte) (*parsedResponse, error) {
	var result receivedCoinMarketCapAPIUahResponse

	err := json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, services.ErrAPIResponseUnmarshallError
	}

	if len(result.Data) == 0 {
		return nil, services.ErrAPIRequestUnsuccessful
	}

	resultUahData := result.Data[0].Quote.Uah
	price := resultUahData.Price

	timestamp, err := time.Parse("2006-01-02T15:04:05.999Z", resultUahData.Time)
	if err != nil {
		return nil, services.ErrAPIResponseUnmarshallError
	}

	return &parsedResponse{
		rate: price,
		time: timestamp,
	}, nil
}
