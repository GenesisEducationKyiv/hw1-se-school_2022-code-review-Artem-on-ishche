package services

import "fmt"

type EmailAddressesStorage interface {
	IsEmailAddressAlreadySaved(emailAddress EmailAddress) bool
	AddEmailAddress(emailAddress EmailAddress) error
	GetEmailAddresses() []string
}

var ErrEmailAddressAlreadyExists = func(emailAddress string) error {
	return fmt.Errorf("email address %v already exists", emailAddress)
}

func AddEmailAddress(storage EmailAddressesStorage, emailAddress EmailAddress) error {
	if storage.IsEmailAddressAlreadySaved(emailAddress) {
		return ErrEmailAddressAlreadyExists(string(emailAddress))
	}

	return storage.AddEmailAddress(emailAddress)
}
