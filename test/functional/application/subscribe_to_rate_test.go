package application

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type inMemoryEmailAddressesRepository struct {
	emailAddresses []models.EmailAddress
}

func newInMemoryEmailAddressesStorage(emailAddresses []models.EmailAddress) inMemoryEmailAddressesRepository {
	return inMemoryEmailAddressesRepository{emailAddresses}
}

func (repository *inMemoryEmailAddressesRepository) IsSaved(emailAddress models.EmailAddress) bool {
	for _, address := range repository.emailAddresses {
		if address == emailAddress {
			return true
		}
	}

	return false
}

func (repository *inMemoryEmailAddressesRepository) Add(emailAddress models.EmailAddress) error {
	repository.emailAddresses = append(repository.emailAddresses, emailAddress)

	return nil
}

func (repository *inMemoryEmailAddressesRepository) GetAll() []string {
	var emailAddressStrings []string

	for _, address := range repository.emailAddresses {
		emailAddressStrings = append(emailAddressStrings, string(address))
	}

	return emailAddressStrings
}

func TestAddNewEmailAddress(t *testing.T) {
	address, _ := models.NewEmailAddress("my_address@domain.extension")
	storage := newInMemoryEmailAddressesStorage([]models.EmailAddress{})
	addEmailAddressImpl := application.NewAddEmailAddressServiceImpl(&storage)

	err := addEmailAddressImpl.AddEmailAddress(*address)

	assert.Nil(t, err)
}

func TestAddExistingEmailAddress(t *testing.T) {
	addressString := "my_address@domain.extension"
	address, _ := models.NewEmailAddress(addressString)
	storage := newInMemoryEmailAddressesStorage([]models.EmailAddress{*address})
	addEmailAddressImpl := application.NewAddEmailAddressServiceImpl(&storage)

	err := addEmailAddressImpl.AddEmailAddress(*address)

	assert.NotNil(t, err)
	assert.Equal(t, services.ErrEmailAddressAlreadyExists(addressString), err)
}

func TestSuccessiveAddingTheSameEmailAddress(t *testing.T) {
	addressString := "my_address@domain.extension"
	address, _ := models.NewEmailAddress(addressString)
	storage := newInMemoryEmailAddressesStorage([]models.EmailAddress{})
	addEmailAddressImpl := application.NewAddEmailAddressServiceImpl(&storage)

	err := addEmailAddressImpl.AddEmailAddress(*address)
	assert.Nil(t, err)

	err = addEmailAddressImpl.AddEmailAddress(*address)
	assert.NotNil(t, err)
	assert.Equal(t, services.ErrEmailAddressAlreadyExists(addressString), err)
}

func TestAddMultipleNewEmailAddresses(t *testing.T) {
	addressStrings := []string{
		"address@what.com",
		"artem.mykytyshyn@gmail.com",
		"someone@some.mail",
	}
	storage := newInMemoryEmailAddressesStorage([]models.EmailAddress{})
	addEmailAddressImpl := application.NewAddEmailAddressServiceImpl(&storage)

	for _, addressString := range addressStrings {
		address, _ := models.NewEmailAddress(addressString)

		err := addEmailAddressImpl.AddEmailAddress(*address)

		assert.Nil(t, err)
	}
}
