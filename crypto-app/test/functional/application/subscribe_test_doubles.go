package application

import (
	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
	"gses2.app/api/test/functional/publicmocks"
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

func (repository *inMemoryEmailAddressesRepository) AssociatedCurrencyPair() *models.CurrencyPair {
	return &btcUahPair
}

type inMemoryEmailAddressesRepositoryGetter struct {
	repo inMemoryEmailAddressesRepository
}

func newInMemoryEmailAddressesRepositoryGetter(repo inMemoryEmailAddressesRepository) *inMemoryEmailAddressesRepositoryGetter {
	return &inMemoryEmailAddressesRepositoryGetter{repo: repo}
}

func (g inMemoryEmailAddressesRepositoryGetter) GetEmailAddressesRepositories(
	currencyPair *models.CurrencyPair,
) []services.EmailAddressesRepository {
	repos := make([]services.EmailAddressesRepository, 0)
	if *currencyPair == btcUahPair {
		return repos
	}

	repos = append(repos, &g.repo)

	return repos
}

func (g inMemoryEmailAddressesRepositoryGetter) GetAllEmailAddressesRepositories() []services.EmailAddressesRepository {
	repos := make([]services.EmailAddressesRepository, 0)
	repos = append(repos, &g.repo)

	return repos
}

func getSubscribeService(addresses []models.EmailAddress) application.RateSubscriptionService {
	storage := newInMemoryEmailAddressesStorage(addresses)
	repoGetter := newInMemoryEmailAddressesRepositoryGetter(storage)
	subscribeServiceImpl := application.NewSubscribeToRateServiceImpl(repoGetter, publicmocks.EmptyCustomersService, publicmocks.EmptyLogger)

	return subscribeServiceImpl
}
