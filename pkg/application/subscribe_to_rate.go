package application

import (
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type AddEmailAddressService interface {
	AddEmailAddress(emailAddress models.EmailAddress) error
}

type addEmailAddressServiceImpl struct {
	repository services.EmailAddressesRepository
}

func NewAddEmailAddressServiceImpl(repository services.EmailAddressesRepository) AddEmailAddressService {
	return &addEmailAddressServiceImpl{repository}
}

func (addEmailService *addEmailAddressServiceImpl) AddEmailAddress(emailAddress models.EmailAddress) error {
	isEmailSaved, err := addEmailService.repository.IsSaved(emailAddress)
	if err != nil {
		return err
	}

	if isEmailSaved {
		return services.ErrEmailAddressAlreadyExists(string(emailAddress))
	}

	return addEmailService.repository.Add(emailAddress)
}

type SubscribeToRateService interface {
	Subscribe(emailAddress *models.EmailAddress, currencyPair *models.CurrencyPair) error
}

type subscribeToRateServiceImpl struct {
	repoGetter services.EmailAddressesRepositoryGetter
}

func NewSubscribeToRateServiceImpl(repoGetter services.EmailAddressesRepositoryGetter) *subscribeToRateServiceImpl {
	return &subscribeToRateServiceImpl{repoGetter: repoGetter}
}

func (s subscribeToRateServiceImpl) Subscribe(emailAddress *models.EmailAddress, currencyPair *models.CurrencyPair) error {
	repository := s.repoGetter.GetEmailAddressesRepository(currencyPair)

	isEmailSaved, err := repository.IsSaved(*emailAddress)
	if err != nil {
		return err
	}

	if isEmailSaved {
		return services.ErrEmailAddressAlreadyExists(emailAddress.String())
	}

	return repository.Add(*emailAddress)
}
