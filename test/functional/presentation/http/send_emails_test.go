package http

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gses2.app/api/pkg/domain/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendEmailsHandlerWhenEverythingIsOk(t *testing.T) {
	setSendBtcToUahRateEmailsTestFunctionToReturnNoError()

	recorder := makeSendEmailsRequest()

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestSendEmailsHandlerWhenApiRequestFailed(t *testing.T) {
	setSendBtcToUahRateEmailsTestFunctionToReturnError(services.ErrAPIRequestUnsuccessful)

	recorder := makeSendEmailsRequest()

	assert.Equal(t, http.StatusBadGateway, recorder.Code)
}

func TestSendEmailsHandlerWhenSomethingElseFailed(t *testing.T) {
	setSendBtcToUahRateEmailsTestFunctionToReturnError(fmt.Errorf("some unknown error"))

	recorder := makeSendEmailsRequest()

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
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

func makeSendEmailsRequest() *httptest.ResponseRecorder {
	router, recorder := getRouterAndRecorder()
	request, _ := http.NewRequest("POST", "/sendEmails", nil)
	router.ServeHTTP(recorder, request)

	return recorder
}
