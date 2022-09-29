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

func TestSubscribeRequestHandlerWhenParameterIsMissing(t *testing.T) {
	setAddEmailAddressFunctionToReturnNoError()

	recorder := makeSubscribeRequest(getURLWithoutRequiredParameter())

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "parameter")
}

func TestSubscribeRequestHandlerWhenEmailParameterIsWrong(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturnRateWithPrice(100)
	setAddEmailAddressFunctionToReturnNoError()

	recorder := makeSubscribeRequest(getURLWithWrongEmailParameter())

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.NotContains(t, recorder.Body.String(), "parameter")
	assert.Contains(t, recorder.Body.String(), "email address")
}

func TestSubscribeRequestHandlerWhenEmailIsAlreadySaved(t *testing.T) {
	setAddEmailAddressFunctionToReturnEmailAlreadyExistsError()

	recorder := makeSubscribeRequest(getCorrectURL())

	assert.Equal(t, http.StatusConflict, recorder.Code)
}

func TestSubscribeRequestHandlerWhenSomeErrorOccurs(t *testing.T) {
	setAddEmailAddressFunctionToReturnUnknownError()

	recorder := makeSubscribeRequest(getCorrectURL())

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestSubscribeRequestHandlerWhenEverythingIsOk(t *testing.T) {
	setAddEmailAddressFunctionToReturnNoError()

	recorder := makeSubscribeRequest(getCorrectURL())

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func setAddEmailAddressFunctionToReturnNoError() {
	addEmailAddressTestFunction = func(emailAddress models.EmailAddress) error {
		return nil
	}
}

func setAddEmailAddressFunctionToReturnEmailAlreadyExistsError() {
	addEmailAddressTestFunction = func(emailAddress models.EmailAddress) error {
		return services.ErrEmailAddressAlreadyExists(string(emailAddress))
	}
}

func setAddEmailAddressFunctionToReturnUnknownError() {
	addEmailAddressTestFunction = func(emailAddress models.EmailAddress) error {
		return fmt.Errorf("some unknown error")
	}
}

func getURLWithoutRequiredParameter() string {
	return "/subscribe"
}

func getURLWithWrongEmailParameter() string {
	return "/subscribe?email=not.a.valid.email.com"
}

func getCorrectURL() string {
	return "/subscribe?email=name@mail.com"
}

func makeSubscribeRequest(url string) *httptest.ResponseRecorder {
	router, recorder := getRouterAndRecorder()
	request, _ := http.NewRequest("POST", url, nil)
	router.ServeHTTP(recorder, request)

	return recorder
}
