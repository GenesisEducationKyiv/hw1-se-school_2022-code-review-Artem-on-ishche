package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/services"
)

type sendEmailsFunction func() error

var sendBtcToUahRateEmailsTestFunction sendEmailsFunction

type sendBtcToUahRateEmailsServiceTestDouble struct{}

func (service sendBtcToUahRateEmailsServiceTestDouble) SendBtcToUahRateEmails() error {
	return sendBtcToUahRateEmailsTestFunction()
}

var testSendEmailsHandler = NewSendEmailsRequestHandler(sendBtcToUahRateEmailsServiceTestDouble{})

func TestSendEmailsHandlerWhenEverythingIsOk(t *testing.T) {
	setSendBtcToUahRateEmailsTestFunctionToReturnNoError()

	response := testSendEmailsHandler.HandleRequest(nil)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestSendEmailsHandlerWhenApiRequestFailed(t *testing.T) {
	setSendBtcToUahRateEmailsTestFunctionToReturnError(services.ErrAPIRequestUnsuccessful)

	response := testSendEmailsHandler.HandleRequest(nil)

	assert.Equal(t, http.StatusBadGateway, response.StatusCode)
}

func TestSendEmailsHandlerWhenSomethingElseFailed(t *testing.T) {
	setSendBtcToUahRateEmailsTestFunctionToReturnError(fmt.Errorf("some unknown error"))

	response := testSendEmailsHandler.HandleRequest(nil)

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func setSendBtcToUahRateEmailsTestFunctionToReturnNoError() {
	sendBtcToUahRateEmailsTestFunction = func() error {
		return nil
	}
}

func setSendBtcToUahRateEmailsTestFunctionToReturnError(err error) {
	sendBtcToUahRateEmailsTestFunction = func() error {
		return err
	}
}
