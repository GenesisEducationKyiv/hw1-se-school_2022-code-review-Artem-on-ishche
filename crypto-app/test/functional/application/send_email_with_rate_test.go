package application

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/test/functional/publicmocks"
)

func TestSendRateEmails_WithSpy_CallsEmailSender(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturnRateWithPrice(0)
	setSendEmailsTestFunctionToCountCalls()

	rateService, storage, sender := getRateServiceRepoGetterAndSenderImplementations([]models.EmailAddress{})
	sendEmailsServiceImpl := application.NewSendRateEmailsServiceImpl(key, rateService, storage, sender, publicmocks.EmptyLogger)

	_ = sendEmailsServiceImpl.SendRateEmails(&btcUahPair, key)

	assert.Equal(t, 1, sendEmailsCallCount)
}

func TestThatEmailBodyContainsBtcToUahRate(t *testing.T) {
	rate := 100.23

	setGetRateWithoutErrorFunctionToReturnRateWithPrice(rate)
	setSendEmailsTestFunctionToSaveEmailBody()

	rateService, storage, sender := getRateServiceRepoGetterAndSenderImplementations([]models.EmailAddress{})
	sendEmailsServiceImpl := application.NewSendRateEmailsServiceImpl(key, rateService, storage, sender, publicmocks.EmptyLogger)

	_ = sendEmailsServiceImpl.SendRateEmails(&btcUahPair, key)

	assert.Contains(t, emailBodyInTest, fmt.Sprintf("%v", rate))
}

func TestThatEmailHasCorrectReceivers(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturnRateWithPrice(0)
	setSendEmailsTestFunctionToSaveReceiverAddressStrings()

	receiverAddresses := getReceiverAddressesForEmailReceiversTest()
	rateService, storage, sender := getRateServiceRepoGetterAndSenderImplementations(receiverAddresses)
	sendEmailsServiceImpl := application.NewSendRateEmailsServiceImpl(key, rateService, storage, sender, publicmocks.EmptyLogger)

	_ = sendEmailsServiceImpl.SendRateEmails(&btcUahPair, key)

	assert.Equal(t, receiverAddresses, actualReceiverAddresses)
}

func TestThatEmailSenderDoesNotReturnErrorWhenEverythingIsSuccessful(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturnRateWithPrice(0)
	setSendEmailsTestFunctionToNotDoAnything()

	rateService, storage, sender := getRateServiceRepoGetterAndSenderImplementations([]models.EmailAddress{})
	sendEmailsServiceImpl := application.NewSendRateEmailsServiceImpl(key, rateService, storage, sender, publicmocks.EmptyLogger)

	err := sendEmailsServiceImpl.SendRateEmails(&btcUahPair, key)

	assert.Nil(t, err)
}

func TestThatEmailSenderHandlesApiErrors(t *testing.T) {
	setGetRateFunctionToReturnError(application.ErrAPIRequestUnsuccessful)
	setSendEmailsTestFunctionToReturnError()

	rateService, storage, sender := getRateServiceRepoGetterAndSenderImplementations([]models.EmailAddress{})
	sendEmailsServiceImpl := application.NewSendRateEmailsServiceImpl(key, rateService, storage, sender, publicmocks.EmptyLogger)

	err := sendEmailsServiceImpl.SendRateEmails(&btcUahPair, key)

	assert.NotNil(t, err)
	assert.Equal(t, application.ErrAPIRequestUnsuccessful, err)
}

func TestThatEmailSenderHandlesEmailSendingErrors(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturnRateWithPrice(0)
	setSendEmailsTestFunctionToReturnError()

	rateService, storage, sender := getRateServiceRepoGetterAndSenderImplementations([]models.EmailAddress{})
	sendEmailsServiceImpl := application.NewSendRateEmailsServiceImpl(key, rateService, storage, sender, publicmocks.EmptyLogger)

	err := sendEmailsServiceImpl.SendRateEmails(&btcUahPair, key)

	assert.NotNil(t, err)
	assert.NotEqual(t, application.ErrAPIRequestUnsuccessful, err)
}
