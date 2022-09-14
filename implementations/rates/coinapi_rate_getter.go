package rates

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/resty.v0"

	"gses2.app/api/config"
	"gses2.app/api/services"
)

type receivedCoinAPIResponse struct {
	Time string  `json:"time"`
	Rate float64 `json:"rate"`
}

type CoinAPIClientFactory struct{}

func (factory CoinAPIClientFactory) CreateRateService() services.ExchangeRateService {
	return &exchangeRateCoinAPIClient{
		apiRequestFormat: "https://rest.coinapi.io/v1/exchangerate/%v/%v",
		apiKeyHeader:     "X-CoinAPI-Key",
		apiKeyValue:      config.CoinAPIKeyValue,
	}
}

type exchangeRateCoinAPIClient struct {
	apiRequestFormat string
	apiKeyHeader     string
	apiKeyValue      string

	next *services.ExchangeRateService
}

func (client *exchangeRateCoinAPIClient) SetNext(service *services.ExchangeRateService) {
	client.next = service
}

func (client *exchangeRateCoinAPIClient) GetExchangeRate(from, to services.Currency) (float64, error) {
	rate, err := client.getExchangeRate(from, to)
	if err != nil && client.next != nil {
		return (*client.next).GetExchangeRate(from, to)
	}

	return rate, err
}

func (client *exchangeRateCoinAPIClient) getExchangeRate(from, to services.Currency) (float64, error) {
	resp, err := client.makeAPIRequest(from, to)
	if !isAPIRequestSuccessful(resp, err) {
		return -1, services.ErrAPIRequestUnsuccessful
	}

	var result receivedCoinAPIResponse

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return -1, err
	}

	return result.Rate, nil
}

func (client *exchangeRateCoinAPIClient) makeAPIRequest(from, to services.Currency) (*resty.Response, error) {
	url := client.getAPIRequestForGivenCurrencies(from, to)

	return resty.R().
		SetHeader(client.apiKeyHeader, client.apiKeyValue).
		Get(url)
}

func (client *exchangeRateCoinAPIClient) getAPIRequestForGivenCurrencies(from, to services.Currency) string {
	return fmt.Sprintf(client.apiRequestFormat, from.Name, to.Name)
}

func isAPIRequestSuccessful(resp *resty.Response, err error) bool {
	return err == nil && resp.StatusCode() == http.StatusOK
}
