package rates

import (
	"encoding/json"
	"fmt"

	"gopkg.in/resty.v0"

	"gses2.app/api/config"
	"gses2.app/api/services"
)

type receivedCoinMarketCapAPIUahResponse struct {
	Data []struct {
		Quote struct {
			Uah struct {
				Price float64 `json:"price"`
			} `json:"UAH"`
		} `json:"quote"`
	} `json:"data"`
}

type CoinMarketCapAPIClientFactory struct{}

func (factory CoinMarketCapAPIClientFactory) CreateRateService() services.ExchangeRateService {
	return &exchangeRateCoinMarketCapAPIClient{
		apiRequestFormat: "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?convert=%v&limit=1&start=1",
		apiKeyHeader:     "X-CMC_PRO_API_KEY",
		apiKeyValue:      config.CoinMarketCapAPIKeyValue,
	}
}

type exchangeRateCoinMarketCapAPIClient struct {
	apiRequestFormat string
	apiKeyHeader     string
	apiKeyValue      string

	next *services.ExchangeRateService
}

func (client *exchangeRateCoinMarketCapAPIClient) SetNext(service *services.ExchangeRateService) {
	client.next = service
}

func (client *exchangeRateCoinMarketCapAPIClient) GetExchangeRate(from, to services.Currency) (float64, error) {
	rate, err := client.getExchangeRate(from, to)
	if err != nil && client.next != nil {
		return (*client.next).GetExchangeRate(from, to)
	}

	return rate, err
}

func (client *exchangeRateCoinMarketCapAPIClient) getExchangeRate(from, to services.Currency) (float64, error) {
	if from.Name != "BTC" {
		return -1, services.ErrAPIRequestUnsuccessful
	}

	resp, err := client.makeAPIRequest(to)
	if !isAPIRequestSuccessful(resp, err) {
		return -1, services.ErrAPIRequestUnsuccessful
	}

	var result receivedCoinMarketCapAPIUahResponse

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return -1, err
	}

	if len(result.Data) == 0 {
		return -1, services.ErrAPIRequestUnsuccessful
	}

	return result.Data[0].Quote.Uah.Price, nil
}

func (client *exchangeRateCoinMarketCapAPIClient) makeAPIRequest(to services.Currency) (*resty.Response, error) {
	url := client.getAPIRequestForGivenCurrencies(to)

	return resty.R().
		SetHeader(client.apiKeyHeader, client.apiKeyValue).
		Get(url)
}

func (client *exchangeRateCoinMarketCapAPIClient) getAPIRequestForGivenCurrencies(to services.Currency) string {
	return fmt.Sprintf(client.apiRequestFormat, to.Name)
}
