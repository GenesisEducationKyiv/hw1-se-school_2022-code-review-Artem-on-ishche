package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/antihax/optional"
	mailslurp "github.com/mailslurp/mailslurp-client-go"
	"github.com/stretchr/testify/assert"

	"gses2.app/api/config"
	"gses2.app/api/implementations/email"
	"gses2.app/api/implementations/repos"
)

func TestSendEmailsRoute(t *testing.T) {
	config.LoadEnv()

	testServer := httptest.NewServer(createMux())
	defer testServer.Close()

	mailSlurpClient, ctx, inbox := createClientContextAndInbox()

	saveNewInboxAddress(t, testServer, inbox.EmailAddress)
	requestToSendEmail(t, testServer)

	assertThatEmailHasBeenDelivered(t, inbox, mailSlurpClient, ctx)
}

func createMux() *http.ServeMux {
	genericExchangeRateService := getGenericExchangeRateService()
	emailAddressesRepository := repos.GetEmailAddressesFileRepository()
	emailSender := email.GetEmailClient()

	subscribeRoute := initSubscribeRoute(&emailAddressesRepository)
	sendEmailsRoute := initSendEmailsRoute(&genericExchangeRateService, &emailAddressesRepository, &emailSender)

	mux := http.NewServeMux()

	mux.HandleFunc("/subscribe", subscribeRoute.processRequest)
	mux.HandleFunc("/sendEmails", sendEmailsRoute.processRequest)

	return mux
}

func createClientContextAndInbox() (*mailslurp.APIClient, *context.Context, *mailslurp.InboxDto) {
	ctx := context.WithValue(
		context.Background(),
		mailslurp.ContextAPIKey,
		mailslurp.APIKey{Key: config.MailSlurpAPIKeyValue},
	)

	mailSlurpConfig := mailslurp.NewConfiguration()
	mailSlurpClient := mailslurp.NewAPIClient(mailSlurpConfig)
	inbox, _, _ := mailSlurpClient.InboxControllerApi.CreateInbox(ctx, nil)

	return mailSlurpClient, &ctx, &inbox
}

func saveNewInboxAddress(t *testing.T, testServer *httptest.Server, emailAddress string) {
	t.Helper()

	requestURL := fmt.Sprintf("%v/subscribe?email=%v", testServer.URL, emailAddress)
	request, _ := http.NewRequest("POST", requestURL, nil)
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func requestToSendEmail(t *testing.T, testServer *httptest.Server) {
	t.Helper()

	requestURL := testServer.URL + "/sendEmails"
	request, _ := http.NewRequest("POST", requestURL, nil)
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func assertThatEmailHasBeenDelivered(t *testing.T,
	inbox *mailslurp.InboxDto, client *mailslurp.APIClient, ctx *context.Context,
) {
	t.Helper()

	waitOpts := &mailslurp.WaitForLatestEmailOpts{
		InboxId:    optional.NewInterface(inbox.Id),
		Timeout:    optional.NewInt64(30000),
		UnreadOnly: optional.NewBool(true),
	}
	receivedEmail, _, err := client.WaitForControllerApi.WaitForLatestEmail(*ctx, waitOpts)

	assert.NoError(t, err)
	assert.Contains(t, *receivedEmail.Body, "Зараз 1 біткоїн коштує")
}
