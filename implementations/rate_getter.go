package implementations

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/resty.v0"

	"gses2.app/api/config"
	"gses2.app/api/services"
)

type receivedAPIResponse struct {
	Time         string  `json:"time"`
	AssetIDBase  string  `json:"asset_id_base"`
	AssetIDQuote string  `json:"asset_id_quote"`
	Rate         float64 `json:"rate"`
}

type exchangeRateAPIClient struct {
	apiRequestFormat string
	apiKeyHeader     string
	apiKeyValue      string
}

func GetExchangeRateCoinAPIClient() services.ExchangeRateService {
	return &exchangeRateAPIClient{
		apiRequestFormat: config.APIRequestFormat,
		apiKeyHeader:     config.APIKeyHeader,
		apiKeyValue:      config.APIKeyValue,
	}
}

func (client *exchangeRateAPIClient) GetExchangeRate(from, to services.Currency) (float64, error) {
	resp, err := client.makeAPIRequest(from, to)
	if !isAPIRequestSuccessful(resp, err) {
		return -1, services.ErrAPIRequestUnsuccessful
	}

	var result receivedAPIResponse

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return -1, err
	}

	return result.Rate, nil
}

func (client *exchangeRateAPIClient) makeAPIRequest(from, to services.Currency) (*resty.Response, error) {
	url := client.getAPIRequestForGivenCurrencies(from, to)

	return resty.R().
		SetHeader(client.apiKeyHeader, client.apiKeyValue).
		Get(url)
}

func (client *exchangeRateAPIClient) getAPIRequestForGivenCurrencies(from, to services.Currency) string {
	return fmt.Sprintf(client.apiRequestFormat, from.Name, to.Name)
}

func isAPIRequestSuccessful(resp *resty.Response, err error) bool {
	return err == nil && resp.StatusCode() == http.StatusOK
}
