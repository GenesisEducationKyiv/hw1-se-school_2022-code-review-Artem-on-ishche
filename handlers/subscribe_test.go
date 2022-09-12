package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/services"
)

type addEmailAddressFunction func(emailAddress services.EmailAddress) error

var addEmailAddressTestFunction addEmailAddressFunction

type addEmailAddressServiceTestDouble struct{}

func (service addEmailAddressServiceTestDouble) AddEmailAddress(emailAddress services.EmailAddress) error {
	return addEmailAddressTestFunction(emailAddress)
}

var testSubscribeRequestHandler = subscribeRequestHandler{addEmailAddressServiceTestDouble{}}

func TestSubscribeRequestHandlerWhenParameterIsMissing(t *testing.T) {
	setAddEmailAddressFunctionToReturnNoError()

	response := testSubscribeRequestHandler.HandleRequest(getHTTPRequestWithoutRequiredParameter())

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, response.Message, "parameter")
}

func TestSubscribeRequestHandlerWhenEmailParameterIsWrong(t *testing.T) {
	setAddEmailAddressFunctionToReturnNoError()

	response := testSubscribeRequestHandler.HandleRequest(getHTTPRequestWithWrongEmailParameter())

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.NotContains(t, response.Message, "parameter")
	assert.Contains(t, response.Message, "email address")
}

func TestSubscribeRequestHandlerWhenEmailIsAlreadySaved(t *testing.T) {
	setAddEmailAddressFunctionToReturnEmailAlreadyExistsError()

	response := testSubscribeRequestHandler.HandleRequest(getCorrectHTTPRequest())

	assert.Equal(t, http.StatusConflict, response.StatusCode)
}

func TestSubscribeRequestHandlerWhenSomeErrorOccurs(t *testing.T) {
	setAddEmailAddressFunctionToReturnUnknownError()

	response := testSubscribeRequestHandler.HandleRequest(getCorrectHTTPRequest())

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestSubscribeRequestHandlerWhenEverythingIsOk(t *testing.T) {
	setAddEmailAddressFunctionToReturnNoError()

	response := testSubscribeRequestHandler.HandleRequest(getCorrectHTTPRequest())

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

func getHTTPRequestWithoutRequiredParameter() *http.Request {
	return getHTTPRequest("/subscribe?misspelled_email=name@mail.com")
}

func getHTTPRequestWithWrongEmailParameter() *http.Request {
	return getHTTPRequest("/subscribe?email=not.a.valid.email.com")
}

func getCorrectHTTPRequest() *http.Request {
	return getHTTPRequest("/subscribe?email=name@mail.com")
}

func getHTTPRequest(rawURL string) *http.Request {
	requestURL, _ := url.Parse(rawURL)

	return &http.Request{
		URL: requestURL,
	}
}
