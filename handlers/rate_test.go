package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/services"
)

type getRateFunction func() (float64, error)

var getRateTestFunction getRateFunction

type btcToUahServiceTestDouble struct{}

func (service btcToUahServiceTestDouble) GetBtcToUahRate() (float64, error) {
	return getRateTestFunction()
}

var testBtcToUahHandler = btcToUahRateRequestHandler{btcToUahServiceTestDouble{}}

func TestBtcToUahRateRequestHandlerWhenApiRequestIsSuccessful(t *testing.T) {
	rate := 100.001
	setGetRateWithoutErrorFunctionToReturn(rate)

	response := testBtcToUahHandler.HandleRequest(nil)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, response.Message, fmt.Sprintf("%v", rate))
}

func TestBtcToUahRateRequestHandlerWhenApiRequestFailed(t *testing.T) {
	setGetRateFunctionToReturnError(services.ErrAPIRequestUnsuccessful)

	response := testBtcToUahHandler.HandleRequest(nil)

	assert.Equal(t, http.StatusBadGateway, response.StatusCode)
}

func TestBtcToUahRateRequestHandlerWhenSomethingElseFailed(t *testing.T) {
	setGetRateFunctionToReturnError(fmt.Errorf("some error"))

	response := testBtcToUahHandler.HandleRequest(nil)

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func setGetRateWithoutErrorFunctionToReturn(rate float64) {
	getRateTestFunction = func() (float64, error) {
		return rate, nil
	}
}

func setGetRateFunctionToReturnError(err error) {
	getRateTestFunction = func() (float64, error) {
		return -1, err
	}
}
