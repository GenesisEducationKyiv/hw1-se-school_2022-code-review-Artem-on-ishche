package rates

import (
	"testing"

	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/infrastructure/rates"
)

func TestCoinMarketCapAPIClient(t *testing.T) {
	config.LoadEnv()

	coinAPIClient := rates.CoinMarketCapAPIClientFactory{}.CreateRateService()

	testRateAPIClient(t, coinAPIClient)
}
