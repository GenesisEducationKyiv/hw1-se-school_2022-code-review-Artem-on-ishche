package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
	httpPresentation "gses2.app/api/pkg/presentation/http"
)

type getRateFunction func() (*models.ExchangeRate, error)

var getRateTestFunction getRateFunction

type btcToUahServiceTestDouble struct{}

func (service btcToUahServiceTestDouble) GetBtcToUahRate() (*models.ExchangeRate, error) {
	return getRateTestFunction()
}

var testBtcToUahHandler = httpPresentation.BtcToUahRateRequestHandler{BtcToUahService: btcToUahServiceTestDouble{}}

func TestBtcToUahRateRequestHandlerWhenApiRequestIsSuccessful(t *testing.T) {
	price := 100.001
	setGetRateWithoutErrorFunctionToReturnRateWithPrice(price)

	response := testBtcToUahHandler.GetResponse(nil)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, response.Message, fmt.Sprintf("%v", price))
}

func TestBtcToUahRateRequestHandlerWhenApiRequestFailed(t *testing.T) {
	setGetRateFunctionToReturnError(services.ErrAPIRequestUnsuccessful)

	response := testBtcToUahHandler.GetResponse(nil)

	assert.Equal(t, http.StatusBadGateway, response.StatusCode)
}

func TestBtcToUahRateRequestHandlerWhenSomethingElseFailed(t *testing.T) {
	setGetRateFunctionToReturnError(fmt.Errorf("some error"))

	response := testBtcToUahHandler.GetResponse(nil)

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func setGetRateWithoutErrorFunctionToReturnRateWithPrice(price float64) {
	getRateTestFunction = func() (*models.ExchangeRate, error) {
		return models.NewExchangeRate(
			models.NewCurrencyPair(models.NewCurrency("btc"), models.NewCurrency("uah")),
			price,
		), nil
	}
}

func setGetRateFunctionToReturnError(err error) {
	getRateTestFunction = func() (*models.ExchangeRate, error) {
		return nil, err
	}
}
