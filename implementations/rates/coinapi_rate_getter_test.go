package rates

import (
	"testing"

	"gses2.app/api/config"
)

func TestCoinAPIClient(t *testing.T) {
	config.LoadEnv()

	coinAPIClient := CoinAPIClientFactory{}.CreateRateService()

	testRateAPIClient(t, coinAPIClient)
}
