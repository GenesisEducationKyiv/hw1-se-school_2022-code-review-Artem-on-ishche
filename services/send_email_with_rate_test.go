package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sendEmailsFunction func(email Email, receiverAddresses []string) error

var sendEmailsTestFunction sendEmailsFunction

type spyEmailSender struct{}

func (sender *spyEmailSender) SendEmails(email Email, receiverAddresses []string) error {
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
)
var actualReceiverAddressStrings []string

func TestThatEmailSenderIsCalled(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturn(0)
	setSendEmailsTestFunctionToCountCalls()

	rateService, storage, sender := getRateServiceStorageAndSenderImplementations([]EmailAddress{})
	sendEmailsServiceImpl := NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	_ = sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.Equal(t, 1, sendEmailsCallCount)
}

func TestThatEmailBodyContainsBtcToUahRate(t *testing.T) {
	rate := 100.23

	setGetRateWithoutErrorFunctionToReturn(rate)
	setSendEmailsTestFunctionToSaveEmailBody()

	rateService, storage, sender := getRateServiceStorageAndSenderImplementations([]EmailAddress{})
	sendEmailsServiceImpl := NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	_ = sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.Contains(t, emailBodyInTest, fmt.Sprintf("%v", rate))
}

func TestThatEmailHasCorrectReceivers(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturn(0)
	setSendEmailsTestFunctionToSaveReceiverAddressStrings()

	receiverAddresses := getReceiverAddressesForEmailReceiversTest()
	rateService, storage, sender := getRateServiceStorageAndSenderImplementations(receiverAddresses)
	sendEmailsServiceImpl := NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	_ = sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.Equal(t, actualReceiverAddressStrings, receiverAddressStringsForEmailReceiversTest)
}

func TestThatEmailSenderDoesNotReturnErrorWhenEverythingIsSuccessful(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturn(0)
	setSendEmailsTestFunctionToNotDoAnything()

	rateService, storage, sender := getRateServiceStorageAndSenderImplementations([]EmailAddress{})
	sendEmailsServiceImpl := NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	err := sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.Nil(t, err)
}

func TestThatEmailSenderHandlesApiErrors(t *testing.T) {
	setGetRateFunctionToReturnError(ErrAPIRequestUnsuccessful)
	setSendEmailsTestFunctionToReturnError()

	rateService, storage, sender := getRateServiceStorageAndSenderImplementations([]EmailAddress{})
	sendEmailsServiceImpl := NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	err := sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.NotNil(t, err)
	assert.Equal(t, ErrAPIRequestUnsuccessful, err)
}

func TestThatEmailSenderHandlesEmailSendingErrors(t *testing.T) {
	setGetRateWithoutErrorFunctionToReturn(0)
	setSendEmailsTestFunctionToReturnError()

	rateService, storage, sender := getRateServiceStorageAndSenderImplementations([]EmailAddress{})
	sendEmailsServiceImpl := NewSendBtcToUahRateEmailsServiceImpl(rateService, storage, sender)

	err := sendEmailsServiceImpl.SendBtcToUahRateEmails()

	assert.NotNil(t, err)
	assert.NotEqual(t, ErrAPIRequestUnsuccessful, err)
}

func getRateServiceStorageAndSenderImplementations(receiverAddresses []EmailAddress) (
	BtcToUahRateService, EmailAddressesRepository, EmailSender,
) {
	rateService := NewBtcToUahServiceImpl(&exchangeRateServiceTestDouble{})
	storage := newInMemoryEmailAddressesStorage(receiverAddresses)
	sender := spyEmailSender{}

	return rateService, &storage, &sender
}

func setGetRateWithoutErrorFunctionToReturn(rate float64) {
	getRateTestFunction = func(from, to Currency) (float64, error) {
		return rate, nil
	}
}

func setGetRateFunctionToReturnError(err error) {
	getRateTestFunction = func(from, to Currency) (float64, error) {
		return -1, err
	}
}

func setSendEmailsTestFunctionToNotDoAnything() {
	sendEmailsTestFunction = func(email Email, receiverAddresses []string) error {
		return nil
	}
}

func setSendEmailsTestFunctionToCountCalls() {
	sendEmailsTestFunction = func(email Email, receiverAddresses []string) error {
		sendEmailsCallCount++

		return nil
	}
}

func setSendEmailsTestFunctionToSaveEmailBody() {
	sendEmailsTestFunction = func(email Email, receiverAddresses []string) error {
		emailBodyInTest = email.Body

		return nil
	}
}

func setSendEmailsTestFunctionToSaveReceiverAddressStrings() {
	sendEmailsTestFunction = func(email Email, receiverAddresses []string) error {
		actualReceiverAddressStrings = receiverAddresses

		return nil
	}
}

func setSendEmailsTestFunctionToReturnError() {
	sendEmailsTestFunction = func(email Email, receiverAddresses []string) error {
		return fmt.Errorf("email has not been sent")
	}
}

func getReceiverAddressesForEmailReceiversTest() []EmailAddress {
	var receiverAddresses []EmailAddress

	for _, addressString := range receiverAddressStringsForEmailReceiversTest {
		address, _ := NewEmailAddress(addressString)
		receiverAddresses = append(receiverAddresses, *address)
	}

	return receiverAddresses
}
