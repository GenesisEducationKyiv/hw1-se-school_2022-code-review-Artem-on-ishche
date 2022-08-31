package rates

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/resty.v0"

	"gses2.app/api/config"
)

type receivedAPIResponse struct {
	Time         string  `json:"time"`
	AssetIDBase  string  `json:"asset_id_base"`
	AssetIDQuote string  `json:"asset_id_quote"`
	Rate         float64 `json:"rate"`
}

func GetBtcToUahRate() (float64, error) {
	return getRate("BTC", "UAH")
}

func getRate(currencyFrom, currencyTo string) (float64, error) {
	resp, err := makeAPIRequest(currencyFrom, currencyTo)
	if !isAPIRequestSuccessful(resp, err) {
		return -1, err
	}

	var result receivedAPIResponse

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return -1, err
	}

	return result.Rate, nil
}

func makeAPIRequest(currencyFrom, currencyTo string) (*resty.Response, error) {
	url := getAPIRequestForGivenCurrencies(currencyFrom, currencyTo)

	return resty.R().
		SetHeader(config.APIKeyHeader, config.APIKeyValue).
		Get(url)
}

func getAPIRequestForGivenCurrencies(currencyFrom, currencyTo string) string {
	return fmt.Sprintf(config.APIRequestFormat, currencyFrom, currencyTo)
}

func isAPIRequestSuccessful(resp *resty.Response, err error) bool {
	return err == nil && resp.StatusCode() == http.StatusOK
}
