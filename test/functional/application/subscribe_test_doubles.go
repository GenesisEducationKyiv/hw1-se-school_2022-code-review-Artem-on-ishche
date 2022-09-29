package application

import (
	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

var key = "key"

type inMemoryEmailAddressesRepository struct {
	emailAddresses []models.EmailAddress
}

func newInMemoryEmailAddressesStorage(emailAddresses []models.EmailAddress) inMemoryEmailAddressesRepository {
	return inMemoryEmailAddressesRepository{emailAddresses}
}

func (repository *inMemoryEmailAddressesRepository) IsSaved(emailAddress models.EmailAddress) (bool, error) {
	for _, address := range repository.emailAddresses {
		if address == emailAddress {
			return true, nil
		}
	}

	return false, nil
}

func (repository *inMemoryEmailAddressesRepository) Add(emailAddress models.EmailAddress) error {
	repository.emailAddresses = append(repository.emailAddresses, emailAddress)

	return nil
}

func (repository *inMemoryEmailAddressesRepository) GetAll() ([]models.EmailAddress, error) {
	return repository.emailAddresses, nil
}

type inMemoryEmailAddressesRepositoryGetter struct {
	repo inMemoryEmailAddressesRepository
}

func newInMemoryEmailAddressesRepositoryGetter(repo inMemoryEmailAddressesRepository) *inMemoryEmailAddressesRepositoryGetter {
	return &inMemoryEmailAddressesRepositoryGetter{repo: repo}
}

func (g inMemoryEmailAddressesRepositoryGetter) GetEmailAddressesRepository(*models.CurrencyPair) services.EmailAddressesRepository {
	return &g.repo
}

func getSubscribeService(addresses []models.EmailAddress) application.RateSubscriptionService {
	storage := newInMemoryEmailAddressesStorage(addresses)
	repoGetter := newInMemoryEmailAddressesRepositoryGetter(storage)
	subscribeServiceImpl := application.NewSubscribeToRateServiceImpl(repoGetter)

	return subscribeServiceImpl
}
