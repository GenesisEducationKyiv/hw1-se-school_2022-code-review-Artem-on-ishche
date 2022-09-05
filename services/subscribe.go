package services

import "fmt"

var ErrEmailAddressAlreadyExists = func(emailAddress string) error {
	return fmt.Errorf("email address %v already exists", emailAddress)
}

type EmailAddressesStorage interface {
	IsEmailAddressAlreadySaved(emailAddress EmailAddress) bool
	AddEmailAddress(emailAddress EmailAddress) error
	GetEmailAddresses() []string
}

type AddEmailAddressService interface {
	AddEmailAddress(emailAddress EmailAddress) error
}

type addEmailAddressServiceImpl struct {
	storage EmailAddressesStorage
}

func NewAddEmailAddressServiceImpl(storage EmailAddressesStorage) AddEmailAddressService {
	return &addEmailAddressServiceImpl{storage}
}

func (addEmailService *addEmailAddressServiceImpl) AddEmailAddress(emailAddress EmailAddress) error {
	if addEmailService.storage.IsEmailAddressAlreadySaved(emailAddress) {
		return ErrEmailAddressAlreadyExists(string(emailAddress))
	}

	return addEmailService.storage.AddEmailAddress(emailAddress)
}
