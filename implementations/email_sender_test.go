package implementations

import (
	"context"
	"testing"

	"github.com/antihax/optional"
	mailslurp "github.com/mailslurp/mailslurp-client-go"
	"github.com/stretchr/testify/assert"

	"gses2.app/api/config"
	"gses2.app/api/services"
)

func TestSendEmailsWithNoReceivers(t *testing.T) {
	email := services.NewEmail("title", "content")

	var receiverAddresses []string

	err := GetEmailClient().SendEmails(*email, receiverAddresses)

	assert.NoError(t, err)
}

func TestSendEmails(t *testing.T) {
	config.LoadEnv()

	client, ctx := createClientAndContext()
	inbox, _, _ := client.InboxControllerApi.CreateInbox(ctx, nil)
	emailToSend := services.NewEmail("title", "content")
	receiverAddresses := []string{inbox.EmailAddress}

	err := GetEmailClient().SendEmails(*emailToSend, receiverAddresses)

	assert.NoError(t, err)
	assertReceivedEmail(t, inbox, *client, ctx)
}

func createClientAndContext() (*mailslurp.APIClient, context.Context) {
	ctx := context.WithValue(
		context.Background(),
		mailslurp.ContextAPIKey,
		mailslurp.APIKey{Key: config.MailSlurpAPIKeyValue},
	)

	mailSlurpConfig := mailslurp.NewConfiguration()

	return mailslurp.NewAPIClient(mailSlurpConfig), ctx
}

func assertReceivedEmail(t *testing.T, inbox mailslurp.InboxDto, client mailslurp.APIClient, ctx context.Context) {
	t.Helper()

	waitOpts := &mailslurp.WaitForLatestEmailOpts{
		InboxId:    optional.NewInterface(inbox.Id),
		Timeout:    optional.NewInt64(30000),
		UnreadOnly: optional.NewBool(true),
	}
	receivedEmail, _, err := client.WaitForControllerApi.WaitForLatestEmail(ctx, waitOpts)

	assert.NoError(t, err)
	assert.Contains(t, *receivedEmail.Body, "content")
}
