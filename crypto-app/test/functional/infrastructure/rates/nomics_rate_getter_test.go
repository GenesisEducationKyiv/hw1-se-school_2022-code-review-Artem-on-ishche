package rates

import (
	"testing"

	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/infrastructure/rates"
	"gses2.app/api/test/functional/publicmocks"
)

func TestNomicsAPIClient(t *testing.T) {
	config.LoadEnv(publicmocks.EmptyLogger)

	coinAPIClient := rates.NomicsAPIClientFactory{}.CreateRateService()

	testRateAPIClient(t, coinAPIClient)
}
