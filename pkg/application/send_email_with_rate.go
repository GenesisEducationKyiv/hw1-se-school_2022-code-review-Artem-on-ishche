package application

import (
	"fmt"

	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type SendBtcToUahRateEmailsService interface {
	SendBtcToUahRateEmails() error
}

type sendBtcToUahRateEmailsServiceImpl struct {
	rateService BtcToUahRateService
	repository  services.EmailAddressesRepository
	sender      services.EmailSender
}

func NewSendBtcToUahRateEmailsServiceImpl(
	rateService BtcToUahRateService, repository services.EmailAddressesRepository, sender services.EmailSender,
) SendBtcToUahRateEmailsService {
	return &sendBtcToUahRateEmailsServiceImpl{
		rateService: rateService,
		repository:  repository,
		sender:      sender,
	}
}

func (sendRateEmailsService *sendBtcToUahRateEmailsServiceImpl) SendBtcToUahRateEmails() error {
	rate, err := sendRateEmailsService.rateService.GetBtcToUahRate()
	if err != nil {
		return err
	}

	email := getEmailWithRate(rate)

	receiverAddresses, err := sendRateEmailsService.repository.GetAll()
	if err != nil {
		return err
	}

	return sendRateEmailsService.sender.SendEmails(email, receiverAddresses)
}

func getEmailWithRate(rate *models.ExchangeRate) models.EmailMessage {
	title := "BTC Quote UAH rate"
	body := fmt.Sprintf("Зараз 1 біткоїн коштує %v грн\n", rate.Price)

	return *models.NewEmail(title, body)
}
