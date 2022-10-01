package services

import (
	"errors"
	"fmt"

	"gses2.app/api/pkg/domain/models"
)

var ErrEmailAddressAlreadyExists = func(emailAddress string) error {
	return fmt.Errorf("email address %v already exists", emailAddress)
}

var ErrEmailStorageFailure = errors.New("failure when working with email storage")

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
