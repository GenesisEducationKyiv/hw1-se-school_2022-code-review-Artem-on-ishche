package main

import (
	"bufio"
	"os"
)

// safelyClose loops while the file is not closed successfully.
func safelyClose(file *os.File) {
	err := file.Close()
	for err != nil {
		err = file.Close()
	}
}

// isEmailAddressSaved checks if the emailAddress is present in the file.
// If an I/O error occurs, it means that the file does not yet exist,
// and isEmailAddressSaved returns false.
func isEmailAddressSaved(emailAddress string) bool {
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

// addEmailAddress appends the emailAddress to the end of the file.
// If the file does not exist, it creates one.
// If any I/O error occurs, it passes it up the call stack.
func addEmailAddress(emailAddress string) error {
	fileName := os.Getenv("FILENAME")

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer safelyClose(file)

	_, err = file.WriteString(emailAddress + "\n")
	return err
}

// getEmailAddresses returns the slice of all email addresses in the file.
// If an I/O error occurs, it means that the file does not yet exist,
// and getEmailAddresses returns an empty slice.
func getEmailAddresses() []string {
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
