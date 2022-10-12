package application

import (
	"fmt"

	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type RateSubscriptionService interface {
	Subscribe(emailAddress *models.EmailAddress, currencyPair *models.CurrencyPair) error
}

type rateSubscriptionServiceImpl struct {
	repoGetter services.EmailAddressesRepositoryGetter
	logger     services.Logger
}

func NewSubscribeToRateServiceImpl(
	repoGetter services.EmailAddressesRepositoryGetter,
	logger services.Logger,
) *rateSubscriptionServiceImpl {
	return &rateSubscriptionServiceImpl{repoGetter: repoGetter, logger: logger}
}

func (s rateSubscriptionServiceImpl) Subscribe(emailAddress *models.EmailAddress, currencyPair *models.CurrencyPair) error {
	s.logger.Debug(fmt.Sprintf("Subscribe() called with emailAddress=%s, currencyPair={%s}",
		emailAddress.String(), currencyPair.String()))

	repository := s.repoGetter.GetEmailAddressesRepositories(currencyPair)[0]
	s.logger.Debug(fmt.Sprintf(
		"repoGetter.GetEmailAddressesRepositories() returned a repository with associated currency pair = {%s}",
		repository.AssociatedCurrencyPair().String()))

	isEmailSaved, err := repository.IsSaved(*emailAddress)
	s.logger.Debug(fmt.Sprintf("repository.IsSaved() returned isEmailSaved=%v, err={%s}", isEmailSaved, err))

	if err != nil {
		return err
	} else if isEmailSaved {
		return ErrEmailAddressAlreadyExists(emailAddress.String())
	}

	err = repository.Add(*emailAddress)
	s.logger.Debug(fmt.Sprintf("repository.Add() returned err={%v}", err))

	return nil
}
