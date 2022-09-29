package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

func TestBtcToUahRateRequestHandlerWhenApiRequestIsSuccessful(t *testing.T) {
	price := 100.001
	setGetRateWithoutErrorFunctionToReturnRateWithPrice(price)

	recorder := makeRateRequest()

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), fmt.Sprintf("%v", price))
}

func TestBtcToUahRateRequestHandlerWhenApiRequestFailed(t *testing.T) {
	setGetRateFunctionToReturnError(services.ErrAPIRequestUnsuccessful)

	recorder := makeRateRequest()

	assert.Equal(t, http.StatusBadGateway, recorder.Code)
}

func TestBtcToUahRateRequestHandlerWhenSomethingElseFailed(t *testing.T) {
	setGetRateFunctionToReturnError(fmt.Errorf("some error"))

	recorder := makeRateRequest()

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
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

func makeRateRequest() *httptest.ResponseRecorder {
	router, recorder := getRouterAndRecorder()
	request, _ := http.NewRequest("GET", "/rate", nil)
	router.ServeHTTP(recorder, request)

	return recorder
}
