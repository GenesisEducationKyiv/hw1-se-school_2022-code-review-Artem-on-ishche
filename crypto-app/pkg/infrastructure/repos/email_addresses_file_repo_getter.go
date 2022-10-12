package repos

import (
	"fmt"
	"os"
	"strings"

	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

const (
	dirName       = "data/"
	fileExtension = ".txt"
)

type EmailAddressesFileRepoGetter struct {
	repos  map[string]services.EmailAddressesRepository
	logger services.Logger
}

func NewEmailAddressesFileRepoGetter(logger services.Logger) *EmailAddressesFileRepoGetter {
	return &EmailAddressesFileRepoGetter{
		repos:  make(map[string]services.EmailAddressesRepository),
		logger: logger,
	}
}

func (getter EmailAddressesFileRepoGetter) GetEmailAddressesRepositories(
	currencyPair *models.CurrencyPair,
) []services.EmailAddressesRepository {
	getter.logger.Debug("GetEmailAddressesRepositories() called with currencyPair = " + currencyPair.String())

	getter.loadAllRepos()

	filteredRepos := make([]services.EmailAddressesRepository, 0)

	if currencyPair == nil {
		return filteredRepos
	}

	for pairString, repo := range getter.repos {
		if satisfiesRequest(pairString, currencyPair) {
			filteredRepos = append(filteredRepos, repo)
		}
	}

	getter.logger.Debug(fmt.Sprintf("fileteredRepos after filtering: %v", filteredRepos))

	if len(filteredRepos) != 0 {
		return filteredRepos
	}

	getter.logger.Debug("No repo left after filtering, so create a new one for this pair")

	pairString := currencyPair.String()
	repo := NewEmailAddressesFileRepository(dirName, pairString+fileExtension)

	getter.repos[pairString] = repo
	filteredRepos = append(filteredRepos, repo)

	return filteredRepos
}

func (getter EmailAddressesFileRepoGetter) GetAllEmailAddressesRepositories() []services.EmailAddressesRepository {
	getter.logger.Debug("GetAllEmailAddressesRepositories() called")

	getter.loadAllRepos()

	repos := make([]services.EmailAddressesRepository, 0)

	for _, repo := range getter.repos {
		repos = append(repos, repo)
	}

	return repos
}

func (getter EmailAddressesFileRepoGetter) loadAllRepos() {
	getter.logger.Debug("loadAllRepos() called")

	files, err := getter.readOrCreateDir()
	if err != nil || len(getter.repos) >= len(files) {
		return
	}

	for _, file := range files {
		fileName := file.Name()
		pairString := strings.Replace(fileName, fileExtension, "", 1)

		repo := NewEmailAddressesFileRepository(dirName, fileName)
		getter.repos[pairString] = repo
	}
}

func (getter EmailAddressesFileRepoGetter) readOrCreateDir() ([]os.DirEntry, error) {
	files, err := os.ReadDir(dirName)
	if err != nil {
		getter.logger.Error("storage directory does not yet exist")

		_ = os.Mkdir(dirName, os.ModePerm)
	}

	return files, err
}

func satisfiesRequest(pairString string, pair *models.CurrencyPair) bool {
	if !pair.Base.IsEmpty() && !pair.Quote.IsEmpty() {
		return pairString == pair.String()
	}

	currencyStrings := strings.Split(pairString, models.CurrencySeparator)
	if !pair.Base.IsEmpty() {
		return currencyStrings[0] == pair.Base.Name
	} else if !pair.Quote.IsEmpty() {
		return currencyStrings[1] == pair.Quote.Name
	}

	return true
}
