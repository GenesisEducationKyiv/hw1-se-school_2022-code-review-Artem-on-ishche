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
	logger           services.Logger
}

func NewSendRateEmailsServiceImpl(
	validationKey string,
	rateService services.ExchangeRateService,
	repositoryGetter services.EmailAddressesRepositoryGetter,
	sender services.EmailSender,
	logger services.Logger,
) SendRateEmailsService {
	return &sendRateEmailsServiceImpl{
		validationKey:    validationKey,
		rateService:      rateService,
		repositoryGetter: repositoryGetter,
		sender:           sender,
		logger:           logger,
	}
}

func (service *sendRateEmailsServiceImpl) SendRateEmails(passedPair *models.CurrencyPair, key string) error {
	service.logger.Debug(fmt.Sprintf("SendRateEmails() called with passedPair={%s}, key=%s", passedPair.String(), key))

	if key != service.validationKey {
		service.logger.Error(fmt.Sprintf("Provided validation key (%s) doesn't match the correct one (%s)", key, service.validationKey))

		return ErrValidationError
	}

	repos := service.repositoryGetter.GetEmailAddressesRepositories(passedPair)
	service.logger.Debug(fmt.Sprintf("GetEmailAddressesRepositories() returned repos=%v", repos))

	for _, repo := range repos {
		if err := service.sendEmailsForOneRepo(repo); err != nil {
			return err
		}
	}

	return nil
}

func (service *sendRateEmailsServiceImpl) sendEmailsForOneRepo(repo services.EmailAddressesRepository) error {
	repoPair := repo.AssociatedCurrencyPair()

	service.logger.Debug("Sending emails for a repo associated with" + repoPair.String())

	rate, err := service.rateService.GetExchangeRate(*repoPair)
	service.logger.Debug(fmt.Sprintf("rateService.GetExchangeRate() returned rate={%s}, err={%v}", rate.String(), err))

	if err != nil {
		return err
	}

	email := getEmailWithRate(rate)

	repoAddresses, err := repo.GetAll()
	service.logger.Debug(fmt.Sprintf("repo.GetAll() returned err={%v},\nrepoAddresses=%v", err, repoAddresses))

	if err != nil {
		return err
	}

	err = service.sender.SendEmails(*email, repoAddresses)
	service.logger.Debug(fmt.Sprintf("sender.SendEmails() returned err={%v}", err))

	return err
}

func getEmailWithRate(rate *models.ExchangeRate) *models.EmailMessage {
	title := fmt.Sprintf("%s exchange rate", rate.String())
	body := fmt.Sprintf(
		"At the time of %s, 1 %s costs %v %s",
		rate.Timestamp.String(),
		rate.CurrencyPair.Base,
		rate.Price,
		rate.CurrencyPair.Quote,
	)

	return models.NewEmail(title, body)
}
