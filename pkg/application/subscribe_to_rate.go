package application

import (
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type RateSubscriptionService interface {
	Subscribe(emailAddress *models.EmailAddress, currencyPair *models.CurrencyPair) error
}

type rateSubscriptionServiceImpl struct {
	repoGetter services.EmailAddressesRepositoryGetter
}

func NewSubscribeToRateServiceImpl(repoGetter services.EmailAddressesRepositoryGetter) *rateSubscriptionServiceImpl {
	return &rateSubscriptionServiceImpl{repoGetter: repoGetter}
}

func (s rateSubscriptionServiceImpl) Subscribe(emailAddress *models.EmailAddress, currencyPair *models.CurrencyPair) error {
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
