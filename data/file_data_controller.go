package data

import (
	"bufio"
	"os"
)

// IsEmailAddressSaved checks if the emailAddress is present in the file.
// If an I/O error occurs, it means that the file does not yet exist,
// and IsEmailAddressSaved returns false.
func IsEmailAddressSaved(emailAddress string) bool {
	fileName := os.Getenv("FILENAME")

	file, err := os.Open(fileName)
	if err != nil {
		return false
	}
	defer safelyClose(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if scanner.Text() == emailAddress {
			return true
		}
	}

	return false
}

// AddEmailAddress appends the emailAddress to the end of the file.
// If the file does not exist, it creates one.
// If any I/O error occurs, it passes it up the call stack.
func AddEmailAddress(emailAddress string) error {
	fileName := os.Getenv("FILENAME")

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer safelyClose(file)

	_, err = file.WriteString(emailAddress + "\n")
	return err
}

// GetEmailAddresses returns the slice of all email addresses in the file.
// If an I/O error occurs, it means that the file does not yet exist,
// and GetEmailAddresses returns an empty slice.
func GetEmailAddresses() []string {
	fileName := os.Getenv("FILENAME")

	var emailAddresses []string

	file, err := os.Open(fileName)
	if err != nil {
		return emailAddresses
	}
	defer safelyClose(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		emailAddresses = append(emailAddresses, scanner.Text())
	}

	return emailAddresses
}
