package rates

import (
	"encoding/json"
	"fmt"
	"strconv"

	"gopkg.in/resty.v0"

	"gses2.app/api/config"
	"gses2.app/api/services"
)

type receivedNomicsAPIResponse struct {
	Price          string `json:"price"`
	PriceDate      string `json:"price_date"`
	PriceTimestamp string `json:"price_timestamp"`
}

type NomicsAPIClientFactory struct{}

func (factory NomicsAPIClientFactory) CreateRateService() services.ExchangeRateService {
	return &exchangeRateNomicsAPIClient{
		apiRequestFormat: "https://api.nomics.com/v1/currencies/ticker?key=%v&ids=%v&interval=1d&convert=%v",
		apiKeyValue:      config.NomicsAPIKeyValue,
	}
}

type exchangeRateNomicsAPIClient struct {
	apiRequestFormat string
	apiKeyValue      string
}

func (client *exchangeRateNomicsAPIClient) GetExchangeRate(from, to services.Currency) (float64, error) {
	resp, err := client.makeAPIRequest(from, to)
	if !isAPIRequestSuccessful(resp, err) {
		return -1, services.ErrAPIRequestUnsuccessful
	}

	var results []receivedNomicsAPIResponse

	err = json.Unmarshal(resp.Body, &results)
	if err != nil {
		return -1, err
	}

	if len(results) == 0 {
		return -1, services.ErrAPIRequestUnsuccessful
	}

	bitSize := 64

	price, err := strconv.ParseFloat(results[0].Price, bitSize)
	if err != nil {
		return -1, err
	}

	return price, nil
}

func (client *exchangeRateNomicsAPIClient) makeAPIRequest(from, to services.Currency) (*resty.Response, error) {
	url := client.getAPIRequestForGivenCurrencies(from, to)

	return resty.R().Get(url)
}

func (client *exchangeRateNomicsAPIClient) getAPIRequestForGivenCurrencies(from, to services.Currency) string {
	return fmt.Sprintf(client.apiRequestFormat, client.apiKeyValue, from.Name, to.Name)
}
