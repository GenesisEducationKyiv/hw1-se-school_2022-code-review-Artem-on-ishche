package rates

import (
	"testing"

	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/infrastructure/rates"
)

func TestCoinAPIClient(t *testing.T) {
	config.LoadEnv()

	coinAPIClient := rates.CoinAPIClientFactory{}.CreateRateService()

	testRateAPIClient(t, coinAPIClient)
}
