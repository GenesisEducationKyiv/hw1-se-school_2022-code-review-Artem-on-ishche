package application

import (
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type AddEmailAddressService interface {
	AddEmailAddress(emailAddress models.EmailAddress) error
}

type addEmailAddressServiceImpl struct {
	repository services.EmailAddressesRepository
}

func NewAddEmailAddressServiceImpl(repository services.EmailAddressesRepository) AddEmailAddressService {
	return &addEmailAddressServiceImpl{repository}
}

func (addEmailService *addEmailAddressServiceImpl) AddEmailAddress(emailAddress models.EmailAddress) error {
	isEmailSaved, err := addEmailService.repository.IsSaved(emailAddress)
	if err != nil {
		return err
	}

	if isEmailSaved {
		return services.ErrEmailAddressAlreadyExists(string(emailAddress))
	}

	return addEmailService.repository.Add(emailAddress)
}
