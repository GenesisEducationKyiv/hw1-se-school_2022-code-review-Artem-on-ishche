package application

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type sendEmailsFunction func(email models.EmailMessage, receiverAddresses []models.EmailAddress) error

var sendEmailsTestFunction sendEmailsFunction

type spyEmailSender struct{}

func (sender *spyEmailSender) SendEmails(email models.EmailMessage, receiverAddresses []models.EmailAddress) error {
	return sendEmailsTestFunction(email, receiverAddresses)
}

var (
	sendEmailsCallCount                         = 0
	emailBodyInTest                             string
	receiverAddressStringsForEmailReceiversTest = []string{
		"artem.mykytyshyn@gmail.com",
		"artem.mykytyshyn@ukma.edu.ua",
		"some.other.email.address@email.provider",
	}
	actualReceiverAddresses []models.EmailAddress
)

func TestThatEmailSenderIsCalled(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturnRateWithPrice(0)
	setSendEmailsTestFunctionToCountCalls()

	rateService, storage, sender := getRateServiceStorageAndSenderImplementations([]models.EmailAddress{})
	sendEmailsServiceImpl := application.NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	_ = sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.Equal(t, 1, sendEmailsCallCount)
}

func TestThatEmailBodyContainsBtcToUahRate(t *testing.T) {
	rate := 100.23

	setGetRateWithoutErrorFunctionToReturnRateWithPrice(rate)
	setSendEmailsTestFunctionToSaveEmailBody()

	rateService, storage, sender := getRateServiceStorageAndSenderImplementations([]models.EmailAddress{})
	sendEmailsServiceImpl := application.NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	_ = sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.Contains(t, emailBodyInTest, fmt.Sprintf("%v", rate))
}

func TestThatEmailHasCorrectReceivers(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturnRateWithPrice(0)
	setSendEmailsTestFunctionToSaveReceiverAddressStrings()

	receiverAddresses := getReceiverAddressesForEmailReceiversTest()
	rateService, storage, sender := getRateServiceStorageAndSenderImplementations(receiverAddresses)
	sendEmailsServiceImpl := application.NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	_ = sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.Equal(t, receiverAddresses, actualReceiverAddresses)
}

func TestThatEmailSenderDoesNotReturnErrorWhenEverythingIsSuccessful(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturnRateWithPrice(0)
	setSendEmailsTestFunctionToNotDoAnything()

	rateService, storage, sender := getRateServiceStorageAndSenderImplementations([]models.EmailAddress{})
	sendEmailsServiceImpl := application.NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	err := sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.Nil(t, err)
}

func TestThatEmailSenderHandlesApiErrors(t *testing.T) {
	setGetRateFunctionToReturnError(services.ErrAPIRequestUnsuccessful)
	setSendEmailsTestFunctionToReturnError()

	rateService, storage, sender := getRateServiceStorageAndSenderImplementations([]models.EmailAddress{})
	sendEmailsServiceImpl := application.NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	err := sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.NotNil(t, err)
	assert.Equal(t, services.ErrAPIRequestUnsuccessful, err)
}

func TestThatEmailSenderHandlesEmailSendingErrors(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturnRateWithPrice(0)
	setSendEmailsTestFunctionToReturnError()

	rateService, storage, sender := getRateServiceStorageAndSenderImplementations([]models.EmailAddress{})
	sendEmailsServiceImpl := application.NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	err := sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.NotNil(t, err)
	assert.NotEqual(t, services.ErrAPIRequestUnsuccessful, err)
}

func getRateServiceStorageAndSenderImplementations(receiverAddresses []models.EmailAddress) (
	application.BtcToUahRateService, services.EmailAddressesRepository, services.EmailSender,
) {
	rateService := application.NewBtcToUahServiceImpl(&exchangeRateServiceTestDouble{})
	storage := newInMemoryEmailAddressesStorage(receiverAddresses)
	sender := spyEmailSender{}

	return rateService, &storage, &sender
}

func setGetRateWithoutErrorFunctionToReturnRateWithPrice(price float64) {
	getRateTestFunction = func(pair models.CurrencyPair) (*models.ExchangeRate, error) {
		return models.NewExchangeRate(pair, price), nil
	}
}

func setGetRateFunctionToReturnError(err error) {
	getRateTestFunction = func(pair models.CurrencyPair) (*models.ExchangeRate, error) {
		return nil, err
	}
}

func setSendEmailsTestFunctionToNotDoAnything() {
	sendEmailsTestFunction = func(email models.EmailMessage, receiverAddresses []models.EmailAddress) error {
		return nil
	}
}

func setSendEmailsTestFunctionToCountCalls() {
	sendEmailsTestFunction = func(email models.EmailMessage, receiverAddresses []models.EmailAddress) error {
		sendEmailsCallCount++

		return nil
	}
}

func setSendEmailsTestFunctionToSaveEmailBody() {
	sendEmailsTestFunction = func(email models.EmailMessage, receiverAddresses []models.EmailAddress) error {
		emailBodyInTest = email.Body

		return nil
	}
}

func setSendEmailsTestFunctionToSaveReceiverAddressStrings() {
	sendEmailsTestFunction = func(email models.EmailMessage, receiverAddresses []models.EmailAddress) error {
		actualReceiverAddresses = receiverAddresses

		return nil
	}
}

func setSendEmailsTestFunctionToReturnError() {
	sendEmailsTestFunction = func(email models.EmailMessage, receiverAddresses []models.EmailAddress) error {
		return fmt.Errorf("email has not been sent")
	}
}

func getReceiverAddressesForEmailReceiversTest() []models.EmailAddress {
	var receiverAddresses []models.EmailAddress

	for _, addressString := range receiverAddressStringsForEmailReceiversTest {
		address, _ := models.NewEmailAddress(addressString)
		receiverAddresses = append(receiverAddresses, *address)
	}

	return receiverAddresses
}
