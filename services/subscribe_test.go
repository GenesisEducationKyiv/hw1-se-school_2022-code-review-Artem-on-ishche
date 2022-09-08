package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type inMemoryEmailAddressesStorage struct {
	emailAddresses []EmailAddress
}

func newInMemoryEmailAddressesStorage(emailAddresses []EmailAddress) inMemoryEmailAddressesStorage {
	return inMemoryEmailAddressesStorage{emailAddresses}
}

func (storage *inMemoryEmailAddressesStorage) IsEmailAddressAlreadySaved(emailAddress EmailAddress) bool {
	for _, address := range storage.emailAddresses {
		if address == emailAddress {
			return true
		}
	}

	return false
}

func (storage *inMemoryEmailAddressesStorage) AddEmailAddress(emailAddress EmailAddress) error {
	storage.emailAddresses = append(storage.emailAddresses, emailAddress)

	return nil
}

func (storage *inMemoryEmailAddressesStorage) GetEmailAddresses() []string {
	var emailAddressStrings []string

	for _, address := range storage.emailAddresses {
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
