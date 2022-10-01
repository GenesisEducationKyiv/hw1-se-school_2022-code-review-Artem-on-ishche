package repos

import (
	"bufio"
	"log"
	"os"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
)

func safelyClose(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatal("Problem with closing a data file.")
	}
}

func doesFileContainEmailAddress(scanner *bufio.Scanner, emailAddress string) bool {
	for scanner.Scan() {
		if scanner.Text() == emailAddress {
			return true
		}
	}

	return false
}

func scanAddressesFromFile(file *os.File) ([]models.EmailAddress, error) {
	var emailAddresses []models.EmailAddress

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		emailAddress, err := models.NewEmailAddress(text)
		if err != nil {
			return emailAddresses, application.ErrEmailStorageFailure
		}

		emailAddresses = append(emailAddresses, *emailAddress)
	}

	return emailAddresses, nil
}
