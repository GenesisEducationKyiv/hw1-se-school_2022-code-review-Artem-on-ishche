package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type inMemoryEmailAddressesRepository struct {
	emailAddresses []EmailAddress
}

func newInMemoryEmailAddressesStorage(emailAddresses []EmailAddress) inMemoryEmailAddressesRepository {
	return inMemoryEmailAddressesRepository{emailAddresses}
}

func (repository *inMemoryEmailAddressesRepository) IsSaved(emailAddress EmailAddress) bool {
	for _, address := range repository.emailAddresses {
		if address == emailAddress {
			return true
		}
	}

	return false
}

func (repository *inMemoryEmailAddressesRepository) Add(emailAddress EmailAddress) error {
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
	address, _ := NewEmailAddress("my_address@domain.extension")
	storage := newInMemoryEmailAddressesStorage([]EmailAddress{})
	addEmailAddressImpl := NewAddEmailAddressServiceImpl(&storage)

	err := addEmailAddressImpl.AddEmailAddress(*address)

	assert.Nil(t, err)
}

func TestAddExistingEmailAddress(t *testing.T) {
	addressString := "my_address@domain.extension"
	address, _ := NewEmailAddress(addressString)
	storage := newInMemoryEmailAddressesStorage([]EmailAddress{*address})
	addEmailAddressImpl := NewAddEmailAddressServiceImpl(&storage)

	err := addEmailAddressImpl.AddEmailAddress(*address)

	assert.NotNil(t, err)
	assert.Equal(t, ErrEmailAddressAlreadyExists(addressString), err)
}

func TestSuccessiveAddingTheSameEmailAddress(t *testing.T) {
	addressString := "my_address@domain.extension"
	address, _ := NewEmailAddress(addressString)
	storage := newInMemoryEmailAddressesStorage([]EmailAddress{})
	addEmailAddressImpl := NewAddEmailAddressServiceImpl(&storage)

	err := addEmailAddressImpl.AddEmailAddress(*address)
	assert.Nil(t, err)

	err = addEmailAddressImpl.AddEmailAddress(*address)
	assert.NotNil(t, err)
	assert.Equal(t, ErrEmailAddressAlreadyExists(addressString), err)
}

func TestAddMultipleNewEmailAddresses(t *testing.T) {
	addressStrings := []string{
		"address@what.com",
		"artem.mykytyshyn@gmail.com",
		"someone@some.mail",
	}
	storage := newInMemoryEmailAddressesStorage([]EmailAddress{})
	addEmailAddressImpl := NewAddEmailAddressServiceImpl(&storage)

	for _, addressString := range addressStrings {
		address, _ := NewEmailAddress(addressString)

		err := addEmailAddressImpl.AddEmailAddress(*address)

		assert.Nil(t, err)
	}
}
