package rates

import (
	"testing"

	"gses2.app/api/config"
)

func TestCoinMarketCapAPIClient(t *testing.T) {
	config.LoadEnv()

	coinAPIClient := CoinMarketCapAPIClientFactory{}.CreateRateService()

	testRateAPIClient(t, coinAPIClient)
}
