package application

import (
	"errors"
	"fmt"
)

var (
	ErrAPIRequestUnsuccessful     = errors.New("API request has been unsuccessful")
	ErrAPIResponseUnmarshallError = errors.New("error when unmarshalling API response")

	ErrEmailStorageFailure       = errors.New("failure when working with email storage")
	ErrEmailAddressAlreadyExists = func(emailAddress string) error {
		return fmt.Errorf("email address %v already exists", emailAddress)
	}

	ErrCustomersServiceRequestUnsuccessful = errors.New("API request to customers service has been unsuccessful")
)
