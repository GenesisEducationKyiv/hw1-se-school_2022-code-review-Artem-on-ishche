package rates

import (
	"testing"

	"gses2.app/api/config"
)

func TestNomicsAPIClient(t *testing.T) {
	config.LoadEnv()

	coinAPIClient := NomicsAPIClientFactory{}.CreateRateService()

	testRateAPIClient(t, coinAPIClient)
}
