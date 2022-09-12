package services

import "fmt"

var ErrEmailAddressAlreadyExists = func(emailAddress string) error {
	return fmt.Errorf("email address %v already exists", emailAddress)
}

type EmailAddressesRepository interface {
	IsSaved(emailAddress EmailAddress) bool
	Add(emailAddress EmailAddress) error
	GetAll() []string
}

type AddEmailAddressService interface {
	AddEmailAddress(emailAddress EmailAddress) error
}

type addEmailAddressServiceImpl struct {
	repository EmailAddressesRepository
}

func NewAddEmailAddressServiceImpl(repository EmailAddressesRepository) AddEmailAddressService {
	return &addEmailAddressServiceImpl{repository}
}

func (addEmailService *addEmailAddressServiceImpl) AddEmailAddress(emailAddress EmailAddress) error {
	if addEmailService.repository.IsSaved(emailAddress) {
		return ErrEmailAddressAlreadyExists(string(emailAddress))
	}

	return addEmailService.repository.Add(emailAddress)
}
