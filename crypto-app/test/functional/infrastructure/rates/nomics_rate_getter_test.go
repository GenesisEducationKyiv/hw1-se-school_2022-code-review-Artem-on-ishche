package rates

import (
	"testing"

	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/infrastructure/rates"
)

func TestNomicsAPIClient(t *testing.T) {
	config.LoadEnv()

	coinAPIClient := rates.NomicsAPIClientFactory{}.CreateRateService()

	testRateAPIClient(t, coinAPIClient)
}
