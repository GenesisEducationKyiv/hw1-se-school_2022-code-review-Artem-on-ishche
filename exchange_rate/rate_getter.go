package exchange_rate

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v0"
	"net/http"
	"os"
)

// Response is a helper struct to unpack api response values.
type Response struct {
	Time         string  `json:"time"`
	AssetIdBase  string  `json:"asset_id_base"`
	AssetIdQuote string  `json:"asset_id_quote"`
	Rate         float64 `json:"rate"`
}

// getRate is a general function that makes a call to an API asking
// for an exchange rate between two currencies provided as arguments.
// Returns an exchange rate as a float64 value or an error, if any.
func getRate(currencyFrom, currencyTo string) (float64, error) {
	apiKey := os.Getenv("API_KEY")
	url := fmt.Sprintf("https://rest.coinapi.io/v1/exchangerate/%v/%v", currencyFrom, currencyTo)

	resp, err := resty.R().
		SetHeader("X-CoinAPI-Key", apiKey).
		Get(url)

	if err != nil || resp.StatusCode() != http.StatusOK {
		return -1, err
	}

	var result Response
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return -1, err
	}

	return result.Rate, nil
}

// GetBtcUahRate returns a Bitcoin to Ukrainian Hryvnia exchange rate,
// using getRate as a helper function.
func GetBtcUahRate() (float64, error) {
	return getRate("BTC", "UAH")
}
