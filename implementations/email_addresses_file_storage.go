package implementations

import (
	"bufio"
	"gses2.app/api/config"
	"gses2.app/api/services"
	"os"
)

type emailAddressesFileStorage struct {
	filename string
}

func GetEmailAddressesFileStorage() services.EmailAddressesStorage {
	return &emailAddressesFileStorage{filename: config.Filename}
}

func (storage emailAddressesFileStorage) IsEmailAddressAlreadySaved(emailAddress services.EmailAddress) bool {
	file, err := os.Open(storage.filename)
	if err != nil {
		return false
	}

	defer safelyClose(file)
	scanner := bufio.NewScanner(file)

	return doesFileContainEmailAddress(scanner, string(emailAddress))
}

func doesFileContainEmailAddress(scanner *bufio.Scanner, emailAddress string) bool {
	for scanner.Scan() {
		if scanner.Text() == emailAddress {
			return true
		}
	}

	return false
}

func (storage emailAddressesFileStorage) AddEmailAddress(emailAddress services.EmailAddress) error {
	const (
		fileModeFlags       = os.O_APPEND | os.O_CREATE | os.O_WRONLY
		fileModePermutation = 0o644
	)

	file, err := os.OpenFile(storage.filename, fileModeFlags, fileModePermutation)
	if err != nil {
		return err
	}

	defer safelyClose(file)

	_, err = file.WriteString(string(emailAddress) + "\n")

	return err
}

func (storage emailAddressesFileStorage) GetEmailAddresses() []string {
	var emailAddresses []string

	file, err := os.Open(config.Filename)
	if err != nil {
		return emailAddresses
	}

	defer safelyClose(file)

	emailAddresses = scanAddressesFromFile(file)

	return emailAddresses
}

func scanAddressesFromFile(file *os.File) []string {
	var emailAddresses []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		emailAddresses = append(emailAddresses, scanner.Text())
	}

	return emailAddresses
}
