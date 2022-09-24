package repos

import (
	"bufio"
	"os"

	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type emailAddressesFileRepository struct {
	filename string
}

func GetEmailAddressesFileRepository() services.EmailAddressesRepository {
	return &emailAddressesFileRepository{filename: config.Filename}
}

func (repository emailAddressesFileRepository) IsSaved(emailAddress models.EmailAddress) bool {
	file, err := os.Open(repository.filename)
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

func (repository emailAddressesFileRepository) Add(emailAddress models.EmailAddress) error {
	const (
		fileModeFlags       = os.O_APPEND | os.O_CREATE | os.O_WRONLY
		fileModePermutation = 0o644
	)

	file, err := os.OpenFile(repository.filename, fileModeFlags, fileModePermutation)
	if err != nil {
		return err
	}

	defer safelyClose(file)

	_, err = file.WriteString(string(emailAddress) + "\n")

	return err
}

func (repository emailAddressesFileRepository) GetAll() []string {
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
