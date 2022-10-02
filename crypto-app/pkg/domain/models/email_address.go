package models

import (
	"fmt"
	"regexp"
)

type EmailAddress string

func (addr EmailAddress) String() string {
	return string(addr)
}

var ErrEmailAddressNotValid = func(emailAddress string) error {
	return fmt.Errorf("email address %v is not valid", emailAddress)
}

func NewEmailAddress(emailAddress string) (*EmailAddress, error) {
	err := validateEmailAddress(emailAddress)
	if err != nil {
		return nil, err
	}

	return (*EmailAddress)(&emailAddress), nil
}

func validateEmailAddress(emailAddress string) error {
	if !doesEmailAddressMatchPattern(emailAddress) {
		return ErrEmailAddressNotValid(emailAddress)
	}

	return nil
}

func doesEmailAddressMatchPattern(emailAddress string) bool {
	regex := regexp.MustCompile(`^[^@.][^@]*@[^@.]+\.([^@.][^@]+[^@.]|[^@.])$`)

	return regex.MatchString(emailAddress)
}
