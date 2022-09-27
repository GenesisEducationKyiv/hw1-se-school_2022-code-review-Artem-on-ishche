package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/domain/services"
	httpPresentation "gses2.app/api/pkg/presentation/http"
)

type sendEmailsFunction func() error

var sendBtcToUahRateEmailsTestFunction sendEmailsFunction

type sendBtcToUahRateEmailsServiceTestDouble struct{}

func (service sendBtcToUahRateEmailsServiceTestDouble) SendBtcToUahRateEmails() error {
	return sendBtcToUahRateEmailsTestFunction()
}

var testSendEmailsHandler = httpPresentation.SendEmailsRequestHandler{
	SendBtcToUahRateEmailsService: sendBtcToUahRateEmailsServiceTestDouble{},
}

func TestSendEmailsHandlerWhenEverythingIsOk(t *testing.T) {
	setSendBtcToUahRateEmailsTestFunctionToReturnNoError()

	response := testSendEmailsHandler.GetResponse(nil)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestSendEmailsHandlerWhenApiRequestFailed(t *testing.T) {
	setSendBtcToUahRateEmailsTestFunctionToReturnError(services.ErrAPIRequestUnsuccessful)

	response := testSendEmailsHandler.GetResponse(nil)

	assert.Equal(t, http.StatusBadGateway, response.StatusCode)
}

func TestSendEmailsHandlerWhenSomethingElseFailed(t *testing.T) {
	setSendBtcToUahRateEmailsTestFunctionToReturnError(fmt.Errorf("some unknown error"))

	response := testSendEmailsHandler.GetResponse(nil)

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
