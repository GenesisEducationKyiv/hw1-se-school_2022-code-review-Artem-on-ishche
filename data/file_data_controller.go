package data

import (
	"bufio"
	"os"

	"gses2.app/api/config"
)

func IsEmailAddressSaved(emailAddress string) bool {
	file, err := os.Open(config.Filename)
	if err != nil {
		return false
	}

	defer safelyClose(file)
	scanner := bufio.NewScanner(file)

	return doesFileContainEmailAddress(scanner, emailAddress)
}

func doesFileContainEmailAddress(scanner *bufio.Scanner, emailAddress string) bool {
	for scanner.Scan() {
		if scanner.Text() == emailAddress {
			return true
		}
	}

	return false
}

func AddEmailAddress(emailAddress string) error {
	const (
		fileModeFlags       = os.O_APPEND | os.O_CREATE | os.O_WRONLY
		fileModePermutation = 0o644
	)

	file, err := os.OpenFile(config.Filename, fileModeFlags, fileModePermutation)
	if err != nil {
		return err
	}

	defer safelyClose(file)

	_, err = file.WriteString(emailAddress + "\n")

	return err
}

func GetEmailAddresses() []string {
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
