package application

import (
	"errors"
	"fmt"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

var ErrValidationError = errors.New("validation key is wrong")

type SendRateEmailsService interface {
	SendRateEmails(pair *models.CurrencyPair, key string) error
}

type sendRateEmailsServiceImpl struct {
	validationKey    string
	rateService      services.ExchangeRateService
	repositoryGetter services.EmailAddressesRepositoryGetter
	sender           services.EmailSender
}

func NewSendRateEmailsServiceImpl(
	validationKey string,
	rateService services.ExchangeRateService,
	repositoryGetter services.EmailAddressesRepositoryGetter,
	sender services.EmailSender,
) SendRateEmailsService {
	return &sendRateEmailsServiceImpl{
		validationKey:    validationKey,
		rateService:      rateService,
		repositoryGetter: repositoryGetter,
		sender:           sender,
	}
}

func (service *sendRateEmailsServiceImpl) SendRateEmails(passedPair *models.CurrencyPair, key string) error {
	if key != service.validationKey {
		return ErrValidationError
	}

	repos := service.repositoryGetter.GetEmailAddressesRepositories(passedPair)

	for _, repo := range repos {
		if err := service.sendEmailsForOneRepo(repo); err != nil {
			return err
		}
	}

	return nil
}

func (service *sendRateEmailsServiceImpl) sendEmailsForOneRepo(repo services.EmailAddressesRepository) error {
	repoPair := repo.AssociatedCurrencyPair()

	rate, err := service.rateService.GetExchangeRate(*repoPair)
	if err != nil {
		return err
	}

	email := getEmailWithRate(rate)

	repoAddresses, err := repo.GetAll()
	if err != nil {
		return err
	}

	return service.sender.SendEmails(*email, repoAddresses)
}

func getEmailWithRate(rate *models.ExchangeRate) *models.EmailMessage {
	title := fmt.Sprintf("%s exchange rate", rate.String())
	body := fmt.Sprintf("At the time of %s, 1 %s costs %v %s", rate.Timestamp.String(), rate.CurrencyPair.Base, rate.Price, rate.CurrencyPair.Quote)

	return models.NewEmail(title, body)
}
