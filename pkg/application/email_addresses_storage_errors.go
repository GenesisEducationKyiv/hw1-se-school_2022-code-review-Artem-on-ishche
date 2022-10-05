package application

import (
	"errors"
	"fmt"
)

var (
	ErrEmailStorageFailure       = errors.New("failure when working with email storage")
	ErrEmailAddressAlreadyExists = func(emailAddress string) error {
		return fmt.Errorf("email address %v already exists", emailAddress)
	}
)
