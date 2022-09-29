package repos

import (
	"bufio"
	"os"

	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type emailAddressesFileRepository struct {
	filename string
}

func NewEmailAddressesFileRepository(filename string) services.EmailAddressesRepository {
	return &emailAddressesFileRepository{filename: filename}
}

func (repository emailAddressesFileRepository) IsSaved(emailAddress models.EmailAddress) (bool, error) {
	file, err := os.Open(repository.filename)
	if err != nil {
		return false, nil
	}

	defer safelyClose(file)
	scanner := bufio.NewScanner(file)

	return doesFileContainEmailAddress(scanner, string(emailAddress)), nil
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

func (repository emailAddressesFileRepository) GetAll() ([]models.EmailAddress, error) {
	var emailAddresses []models.EmailAddress

	file, err := os.Open(repository.filename)
	if err != nil {
		return emailAddresses, nil
	}

	defer safelyClose(file)

	emailAddresses, err = scanAddressesFromFile(file)
	if err != nil {
		return emailAddresses, err
	}

	return emailAddresses, nil
}

func scanAddressesFromFile(file *os.File) ([]models.EmailAddress, error) {
	var emailAddresses []models.EmailAddress

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		emailAddress, err := models.NewEmailAddress(text)
		if err != nil {
			return emailAddresses, services.ErrEmailStorageFailure
		}

		emailAddresses = append(emailAddresses, *emailAddress)
	}

	return emailAddresses, nil
}
