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
	receiverAddresses := sendRateEmailsService.repository.GetAll()

	return sendRateEmailsService.sender.SendEmails(email, receiverAddresses)
}

func getEmailWithRate(rate float64) models.Email {
	title := "BTC To UAH rate"
	body := fmt.Sprintf("Зараз 1 біткоїн коштує %v грн\n", rate)

	return *models.NewEmail(title, body)
}
