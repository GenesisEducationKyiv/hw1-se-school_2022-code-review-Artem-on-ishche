package repos

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/infrastructure/repos"
)

var filename = "file.txt"

func TestIsEmailAddressAlreadySavedWithEmptyFile(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()

	isSaved, err := repos.NewEmailAddressesFileRepository(filename).IsSaved(emailAddress)

	assert.NoError(t, err)
	assert.False(t, isSaved)
}

func TestIsEmailAddressAlreadySavedWhenItIsNot(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()

	createNonEmptyFileWithEmail("random_email@gmail.com")

	isSaved, err := repos.NewEmailAddressesFileRepository(filename).IsSaved(emailAddress)

	assert.NoError(t, err)
	assert.False(t, isSaved)

	deleteFile()
}

func TestIsEmailAddressAlreadySavedWhenItIs(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()
	createNonEmptyFileWithEmail(string(emailAddress))

	isSaved, err := repos.NewEmailAddressesFileRepository(filename).IsSaved(emailAddress)

	assert.NoError(t, err)
	assert.True(t, isSaved)

	deleteFile()
}

func TestAddEmailAddressWhenFileDoesNotExist(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()

	err := repos.NewEmailAddressesFileRepository(filename).Add(emailAddress)

	assert.NoError(t, err)

	deleteFile()
}

func TestAddEmailAddressWhenFileDoesNotContainIt(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()

	createNonEmptyFileWithEmail("random_email@gmail.com")

	err := repos.NewEmailAddressesFileRepository(filename).Add(emailAddress)

	assert.NoError(t, err)

	deleteFile()
}

func TestAddEmailAddressWhenFileContainsIt(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()

	createNonEmptyFileWithEmail(string(emailAddress))

	err := repos.NewEmailAddressesFileRepository(filename).Add(emailAddress)

	assert.NoError(t, err)

	deleteFile()
}

func TestGetEmailAddressesWhenFileDoesNotExist(t *testing.T) {
	config.LoadEnv()

	emailAddressStrings, err := repos.NewEmailAddressesFileRepository(filename).GetAll()

	assert.NoError(t, err)
	assert.Empty(t, emailAddressStrings)

	deleteFile()
}

func TestGetEmailAddressesWhenFileContainsOneAddress(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()
	createNonEmptyFileWithEmail(string(emailAddress))

	emailAddresses, err := repos.NewEmailAddressesFileRepository(filename).GetAll()

	assert.NoError(t, err)
	assert.Contains(t, emailAddresses, emailAddress)
	assert.Equal(t, 1, len(emailAddresses))

	deleteFile()
}

func TestGetEmailAddressesWhenFileContainsManyAddresses(t *testing.T) {
	providedEmailAddresses := []models.EmailAddress{
		models.EmailAddress("address0@gmail.com"),
		models.EmailAddress("address1@gmail.com"),
		models.EmailAddress("address2@gmail.com"),
		models.EmailAddress("address3@gmail.com"),
		models.EmailAddress("address4@gmail.com"),
	}

	config.LoadEnv()
	createNonEmptyFileWithManyEmails(providedEmailAddresses)

	actualEmailAddresses, err := repos.NewEmailAddressesFileRepository(filename).GetAll()

	assert.NoError(t, err)
	assert.Equal(t, providedEmailAddresses, actualEmailAddresses)

	deleteFile()
}

func TestSuccessiveAddAndIsSavedCalls(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()
	emailAddressesStorage := repos.NewEmailAddressesFileRepository(filename)

	addErr := emailAddressesStorage.Add(emailAddress)
	contains, isSavedErr := emailAddressesStorage.IsSaved(emailAddress)

	assert.NoError(t, addErr)
	assert.NoError(t, isSavedErr)
	assert.True(t, contains)

	deleteFile()
}

func TestSuccessiveAddAndGetCalls(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()
	emailAddressesStorage := repos.NewEmailAddressesFileRepository(filename)

	addErr := emailAddressesStorage.Add(emailAddress)
	savedAddresses, isSavedErr := emailAddressesStorage.GetAll()

	assert.NoError(t, addErr)
	assert.NoError(t, isSavedErr)
	assert.Contains(t, savedAddresses, emailAddress)

	deleteFile()
}

func TestSuccessiveCallsToAllThreeEmailAddressStorageFunctions(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()
	emailAddressesStorage := repos.NewEmailAddressesFileRepository(filename)

	containsBeforeAdding, isSavedErr1 := emailAddressesStorage.IsSaved(emailAddress)
	savedAddressStringsBeforeAdding, getAllErr1 := emailAddressesStorage.GetAll()
	addingErr := emailAddressesStorage.Add(emailAddress)
	containsAfterAdding, isSavedErr2 := emailAddressesStorage.IsSaved(emailAddress)
	savedAddressStringsAfterAdding, getAllErr2 := emailAddressesStorage.GetAll()

	assert.NoError(t, isSavedErr1)
	assert.False(t, containsBeforeAdding)
	assert.NoError(t, getAllErr1)
	assert.Empty(t, savedAddressStringsBeforeAdding)
	assert.NoError(t, addingErr)
	assert.NoError(t, isSavedErr2)
	assert.True(t, containsAfterAdding)
	assert.NoError(t, getAllErr2)
	assert.NotEmpty(t, savedAddressStringsAfterAdding)
	assert.Contains(t, savedAddressStringsAfterAdding, emailAddress)

	deleteFile()
}

func getEmailAddress() models.EmailAddress {
	emailAddress, _ := models.NewEmailAddress("user@mail.com")

	return *emailAddress
}

func createNonEmptyFileWithEmail(emailAddressString string) {
	file, _ := os.Create(filename)
	_, _ = file.WriteString(emailAddressString)
	_ = file.Close()
}

func createNonEmptyFileWithManyEmails(emailAddresses []models.EmailAddress) {
	file, _ := os.Create(filename)

	for _, emailAddress := range emailAddresses {
		_, _ = file.WriteString(emailAddress.String() + "\n")
	}

	_ = file.Close()
}

func deleteFile() {
	_ = os.Remove(filename)
}
