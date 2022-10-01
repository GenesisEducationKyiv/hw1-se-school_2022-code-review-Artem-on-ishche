package services

import (
	"gses2.app/api/pkg/domain/models"
)

type EmailAddressesRepository interface {
	IsSaved(emailAddress models.EmailAddress) (bool, error)
	Add(emailAddress models.EmailAddress) error
	GetAll() ([]models.EmailAddress, error)
	AssociatedCurrencyPair() *models.CurrencyPair
}

type EmailAddressesRepositoryGetter interface {
	GetEmailAddressesRepositories(currencyPair *models.CurrencyPair) []EmailAddressesRepository
	GetAllEmailAddressesRepositories() []EmailAddressesRepository
}
