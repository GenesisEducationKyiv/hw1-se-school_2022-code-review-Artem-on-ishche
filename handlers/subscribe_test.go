package handlers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gses2.app/api/services"
	"net/http"
	"net/url"
	"testing"
)

type addEmailAddressFunction func(emailAddress services.EmailAddress) error

var addEmailAddressTestFunction addEmailAddressFunction

type addEmailAddressServiceTestDouble struct{}

func (service addEmailAddressServiceTestDouble) AddEmailAddress(emailAddress services.EmailAddress) error {
	return addEmailAddressTestFunction(emailAddress)
}

var testSubscribeRequestHandler = NewSubscribeRequestHandler(addEmailAddressServiceTestDouble{})

func TestSubscribeRequestHandlerWhenParameterIsMissing(t *testing.T) {
	setAddEmailAddressFunctionToReturnNoError()

	response := testSubscribeRequestHandler.HandleRequest(getHttpRequestWithoutRequiredParameter())

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, response.Message, "parameter")
}

func TestSubscribeRequestHandlerWhenEmailParameterIsWrong(t *testing.T) {
	setAddEmailAddressFunctionToReturnNoError()

	response := testSubscribeRequestHandler.HandleRequest(getHttpRequestWithWrongEmailParameter())

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.NotContains(t, response.Message, "parameter")
	assert.Contains(t, response.Message, "email address")
}

func TestSubscribeRequestHandlerWhenEmailIsAlreadySaved(t *testing.T) {
	setAddEmailAddressFunctionToReturnEmailAlreadyExistsError()

	response := testSubscribeRequestHandler.HandleRequest(getCorrectHttpRequest())

	assert.Equal(t, http.StatusConflict, response.StatusCode)
}

func TestSubscribeRequestHandlerWhenSomeErrorOccurs(t *testing.T) {
	setAddEmailAddressFunctionToReturnUnknownError()

	response := testSubscribeRequestHandler.HandleRequest(getCorrectHttpRequest())

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestSubscribeRequestHandlerWhenEverythingIsOk(t *testing.T) {
	setAddEmailAddressFunctionToReturnNoError()

	response := testSubscribeRequestHandler.HandleRequest(getCorrectHttpRequest())

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func setAddEmailAddressFunctionToReturnNoError() {
	addEmailAddressTestFunction = func(emailAddress services.EmailAddress) error {
		return nil
	}
}

func setAddEmailAddressFunctionToReturnEmailAlreadyExistsError() {
	addEmailAddressTestFunction = func(emailAddress services.EmailAddress) error {
		return services.ErrEmailAddressAlreadyExists(string(emailAddress))
	}
}

func setAddEmailAddressFunctionToReturnUnknownError() {
	addEmailAddressTestFunction = func(emailAddress services.EmailAddress) error {
		return fmt.Errorf("some unknown error")
	}
}

func getHttpRequestWithoutRequiredParameter() *http.Request {
	return getHttpRequest("/subscribe?misspelled_email=name@mail.com")
}

func getHttpRequestWithWrongEmailParameter() *http.Request {
	return getHttpRequest("/subscribe?email=not.a.valid.email.com")
}

func getCorrectHttpRequest() *http.Request {
	return getHttpRequest("/subscribe?email=name@mail.com")
}

func getHttpRequest(rawURL string) *http.Request {
	requestURL, _ := url.Parse(rawURL)
	return &http.Request{
		URL: requestURL,
	}
}
