package routes

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	mailslurp "github.com/mailslurp/mailslurp-client-go"
	"github.com/stretchr/testify/assert"
	"gses2.app/api/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendEmailsRoute(t *testing.T) {
	config.LoadEnv()
	testServer := httptest.NewServer(http.HandlerFunc(sendEmailsRoute))
	defer testServer.Close()
	mailSlurpClient, ctx, inbox := createClientContextAndInbox()

	saveNewInboxAddress(t, testServer, inbox.EmailAddress)
	requestToSendEmail(t, testServer)

	assertThatEmailHasBeenDelivered(t, *inbox, *mailSlurpClient, *ctx)
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
	requestUrl := fmt.Sprintf("%v/subscribe?email=%v", testServer.URL, emailAddress)
	request, _ := http.NewRequest("POST", requestUrl, nil)
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func requestToSendEmail(t *testing.T, testServer *httptest.Server) {
	requestUrl := testServer.URL + "/sendEmails"
	request, _ := http.NewRequest("POST", requestUrl, nil)
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func assertThatEmailHasBeenDelivered(t *testing.T, inbox mailslurp.InboxDto, client mailslurp.APIClient, ctx context.Context) {
	waitOpts := &mailslurp.WaitForLatestEmailOpts{
		InboxId:    optional.NewInterface(inbox.Id),
		Timeout:    optional.NewInt64(30000),
		UnreadOnly: optional.NewBool(true),
	}
	receivedEmail, _, err := client.WaitForControllerApi.WaitForLatestEmail(ctx, waitOpts)

	assert.NoError(t, err)
	assert.Contains(t, *receivedEmail.Body, "Зараз 1 біткоїн коштує")
}
