package services

import (
	"fmt"

	"gses2.app/api/pkg/domain/models"
)

var ErrEmailAddressAlreadyExists = func(emailAddress string) error {
	return fmt.Errorf("email address %v already exists", emailAddress)
}

type EmailAddressesRepository interface {
	IsSaved(emailAddress models.EmailAddress) bool
	Add(emailAddress models.EmailAddress) error
	GetAll() []string
}
