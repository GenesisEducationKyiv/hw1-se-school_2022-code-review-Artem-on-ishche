package application

import (
	"fmt"

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

func getRateServiceRepoGetterAndSenderImplementations(receiverAddresses []models.EmailAddress) (
	services.ExchangeRateService, services.EmailAddressesRepositoryGetter, services.EmailSender,
) {
	rateService := rateServiceTestDouble{}
	storage := newInMemoryEmailAddressesStorage(receiverAddresses)
	repoGetter := newInMemoryEmailAddressesRepositoryGetter(storage)
	sender := spyEmailSender{}

	return rateService, repoGetter, &sender
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
