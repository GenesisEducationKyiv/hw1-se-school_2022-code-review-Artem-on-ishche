package repos

import (
	"bufio"
	"os"
	"strings"

	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type emailAddressesFileRepository struct {
	path     string
	filename string
}

func NewEmailAddressesFileRepository(path, filename string) services.EmailAddressesRepository {
	return &emailAddressesFileRepository{path: path, filename: filename}
}

func (repository emailAddressesFileRepository) IsSaved(emailAddress models.EmailAddress) (bool, error) {
	file, err := os.Open(repository.fullFileName())
	if err != nil {
		return false, nil
	}

	defer safelyClose(file)
	scanner := bufio.NewScanner(file)

	return doesFileContainEmailAddress(scanner, string(emailAddress)), nil
}

func (repository emailAddressesFileRepository) Add(emailAddress models.EmailAddress) error {
	const (
		fileModeFlags       = os.O_APPEND | os.O_CREATE | os.O_WRONLY
		fileModePermutation = 0o644
	)

	file, err := os.OpenFile(repository.fullFileName(), fileModeFlags, fileModePermutation)
	if err != nil {
		return err
	}

	defer safelyClose(file)

	_, err = file.WriteString(string(emailAddress) + "\n")

	return err
}

func (repository emailAddressesFileRepository) GetAll() ([]models.EmailAddress, error) {
	var emailAddresses []models.EmailAddress

	file, err := os.Open(repository.fullFileName())
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

func (repository emailAddressesFileRepository) AssociatedCurrencyPair() *models.CurrencyPair {
	fileName := strings.ReplaceAll(repository.filename, fileExtension, "")
	tokens := strings.Split(fileName, models.CurrencySeparator)

	base := models.NewCurrency(tokens[0])
	quote := models.NewCurrency(tokens[1])
	pair := models.NewCurrencyPair(base, quote)

	return &pair
}

func (repository emailAddressesFileRepository) fullFileName() string {
	return repository.path + repository.filename
}
