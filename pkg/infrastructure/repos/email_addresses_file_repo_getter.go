package repos

import (
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type EmailAddressesFileRepoGetter struct {
	repos map[string]services.EmailAddressesRepository
}

func NewEmailAddressesFileRepoGetter() *EmailAddressesFileRepoGetter {
	return &EmailAddressesFileRepoGetter{}
}

func (getter EmailAddressesFileRepoGetter) GetEmailAddressesRepository(
	currencyPair *models.CurrencyPair,
) services.EmailAddressesRepository {
	if getter.repos == nil {
		getter.repos = make(map[string]services.EmailAddressesRepository)
	}

	pairString := currencyPair.String()

	repo, found := getter.repos[pairString]
	if !found {
		repo = NewEmailAddressesFileRepository("data/" + pairString + ".txt")
		getter.repos[pairString] = repo
	}

	return repo
}
