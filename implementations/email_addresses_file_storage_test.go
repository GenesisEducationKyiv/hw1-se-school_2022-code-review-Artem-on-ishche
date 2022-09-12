package implementations

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/config"
	"gses2.app/api/services"
)

func TestIsEmailAddressAlreadySavedWithEmptyFile(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()

	isSaved := GetEmailAddressesFileStorage().IsSaved(emailAddress)

	assert.False(t, isSaved)
}

func TestIsEmailAddressAlreadySavedWhenItIsNot(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()

	createNonEmptyFileWithEmail("random_email@gmail.com")

	isSaved := GetEmailAddressesFileStorage().IsSaved(emailAddress)

	assert.False(t, isSaved)

	deleteFile()
}

func TestIsEmailAddressAlreadySavedWhenItIs(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()
	createNonEmptyFileWithEmail(string(emailAddress))

	isSaved := GetEmailAddressesFileStorage().IsSaved(emailAddress)

	assert.True(t, isSaved)

	deleteFile()
}

func TestAddEmailAddressWhenFileDoesNotExist(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()

	err := GetEmailAddressesFileStorage().Add(emailAddress)

	assert.Nil(t, err)

	deleteFile()
}

func TestAddEmailAddressWhenFileDoesNotContainIt(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()

	createNonEmptyFileWithEmail("random_email@gmail.com")

	err := GetEmailAddressesFileStorage().Add(emailAddress)

	assert.Nil(t, err)

	deleteFile()
}

func TestAddEmailAddressWhenFileContainsIt(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()

	createNonEmptyFileWithEmail(string(emailAddress))

	err := GetEmailAddressesFileStorage().Add(emailAddress)

	assert.Nil(t, err)

	deleteFile()
}

func TestGetEmailAddressesWhenFileDoesNotExist(t *testing.T) {
	config.LoadEnv()

	emailAddressStrings := GetEmailAddressesFileStorage().GetAll()

	assert.Empty(t, emailAddressStrings)

	deleteFile()
}

func TestGetEmailAddressesWhenFileContainsOneAddress(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()
	createNonEmptyFileWithEmail(string(emailAddress))

	emailAddressStrings := GetEmailAddressesFileStorage().GetAll()

	assert.Contains(t, emailAddressStrings, string(emailAddress))
	assert.Equal(t, 1, len(emailAddressStrings))

	deleteFile()
}

func TestGetEmailAddressesWhenFileContainsManyAddresses(t *testing.T) {
	providedEmailAddressStrings := []string{
		"address0@gmail.com",
		"address1@gmail.com",
		"address2@gmail.com",
		"address3@gmail.com",
		"address4@gmail.com",
	}

	config.LoadEnv()
	createNonEmptyFileWithManyEmails(providedEmailAddressStrings)

	actualEmailAddressStrings := GetEmailAddressesFileStorage().GetAll()

	assert.Equal(t, providedEmailAddressStrings, actualEmailAddressStrings)
	assert.Contains(t, actualEmailAddressStrings, "address1@gmail.com")

	deleteFile()
}

func TestSuccessiveAddAndIsSavedCalls(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()
	emailAddressesStorage := GetEmailAddressesFileStorage()

	err := emailAddressesStorage.Add(emailAddress)
	contains := emailAddressesStorage.IsSaved(emailAddress)

	assert.Nil(t, err)
	assert.True(t, contains)

	deleteFile()
}

func TestSuccessiveAddAndGetCalls(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()
	emailAddressesStorage := GetEmailAddressesFileStorage()

	err := emailAddressesStorage.Add(emailAddress)
	savedAddressStrings := emailAddressesStorage.GetAll()

	assert.Nil(t, err)
	assert.Contains(t, savedAddressStrings, string(emailAddress))

	deleteFile()
}

func TestSuccessiveCallsToAllThreeEmailAddressStorageFunctions(t *testing.T) {
	config.LoadEnv()

	emailAddress := getEmailAddress()
	emailAddressesStorage := GetEmailAddressesFileStorage()

	containsBeforeAdding := emailAddressesStorage.IsSaved(emailAddress)
	savedAddressStringsBeforeAdding := emailAddressesStorage.GetAll()
	addingErr := emailAddressesStorage.Add(emailAddress)
	containsAfterAdding := emailAddressesStorage.IsSaved(emailAddress)
	savedAddressStringsAfterAdding := emailAddressesStorage.GetAll()

	assert.False(t, containsBeforeAdding)
	assert.Empty(t, savedAddressStringsBeforeAdding)
	assert.Nil(t, addingErr)
	assert.True(t, containsAfterAdding)
	assert.NotEmpty(t, savedAddressStringsAfterAdding)
	assert.Contains(t, savedAddressStringsAfterAdding, string(emailAddress))

	deleteFile()
}

func getEmailAddress() services.EmailAddress {
	emailAddress, _ := services.NewEmailAddress("user@mail.com")

	return *emailAddress
}

func createNonEmptyFileWithEmail(emailAddressString string) {
	file, _ := os.Create(config.Filename)
	_, _ = file.WriteString(emailAddressString)
	_ = file.Close()
}

func createNonEmptyFileWithManyEmails(emailAddressStrings []string) {
	file, _ := os.Create(config.Filename)

	for _, emailAddressString := range emailAddressStrings {
		_, _ = file.WriteString(emailAddressString + "\n")
	}

	_ = file.Close()
}

func deleteFile() {
	_ = os.Remove(config.Filename)
}
