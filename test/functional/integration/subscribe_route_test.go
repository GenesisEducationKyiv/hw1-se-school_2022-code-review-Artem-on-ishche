package integration

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/infrastructure/repos"
	httpPresentation "gses2.app/api/pkg/presentation/http"
)

func TestThatSubscribeRouteReturnsStatusBadRequestWhenNoEmailIsProvided(t *testing.T) {
	config.LoadEnv()

	testServer := getTestServerWithSubscribeRoute()
	defer testServer.Close()

	request, _ := http.NewRequest("POST", testServer.URL+"/subscribe?mail=me@mail.com", nil)
	client := &http.Client{}
	response, err := client.Do(request)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestThatSubscribeRouteReturnsStatusBadRequestWhenProvidedEmailIsWrong(t *testing.T) {
	config.LoadEnv()

	testServer := getTestServerWithSubscribeRoute()
	defer testServer.Close()

	request, _ := http.NewRequest("POST", testServer.URL+"/subscribe?email=wrong@mail@com", nil)
	client := &http.Client{}
	response, err := client.Do(request)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestThatSubscribeRouteReturnsStatusOKWhenProvidedEmailIsNotYetSaved(t *testing.T) {
	config.LoadEnv()

	testServer := getTestServerWithSubscribeRoute()
	defer testServer.Close()

	request, _ := http.NewRequest("POST", testServer.URL+"/subscribe?email=me@mail.com", nil)
	client := &http.Client{}
	response, err := client.Do(request)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	removeStorageFile(t)
}

func TestThatSubscribeRouteReturnsStatusConflictOnConsecutiveCallsWithTheSameEmail(t *testing.T) {
	config.LoadEnv()

	testServer := getTestServerWithSubscribeRoute()
	defer testServer.Close()

	request, _ := http.NewRequest("POST", testServer.URL+"/subscribe?email=repeating_mail@mail.com", nil)
	client := &http.Client{}
	_, _ = client.Do(request)
	response, err := client.Do(request)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, response.StatusCode)

	removeStorageFile(t)
}

func getTestServerWithSubscribeRoute() *httptest.Server {
	emailAddressesRepository := repos.GetEmailAddressesFileRepository()
	addEmailAddressService := application.NewAddEmailAddressServiceImpl(emailAddressesRepository)
	subscribeHandler := httpPresentation.SubscribeRequestHandler{AddEmailAddressService: addEmailAddressService}
	testServer := httptest.NewServer(http.HandlerFunc(httpPresentation.GetHandlerFunction(subscribeHandler)))

	return testServer
}

func removeStorageFile(t *testing.T) {
	t.Helper()

	err := os.Remove(config.Filename)
	if err != nil {
		t.Error("Error when removing a file")
	}
}