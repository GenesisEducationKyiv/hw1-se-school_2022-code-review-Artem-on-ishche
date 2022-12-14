package email

import (
	"context"
	"testing"

	"github.com/antihax/optional"
	mailslurp "github.com/mailslurp/mailslurp-client-go"
	"github.com/stretchr/testify/assert"

	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/domain/models"
	emailsImpl "gses2.app/api/pkg/infrastructure/email"
	"gses2.app/api/test/functional/publicmocks"
)

func TestSendEmailsWithNoReceivers(t *testing.T) {
	email := models.NewEmail("title", "content")

	var receiverAddresses []models.EmailAddress

	err := emailsImpl.GetEmailClient(publicmocks.EmptyLogger).SendEmails(*email, receiverAddresses)

	assert.NoError(t, err)
}

func TestSendEmails(t *testing.T) {
	config.LoadEnv()

	client, ctx := createClientAndContext()
	inbox, _, _ := client.InboxControllerApi.CreateInbox(ctx, nil)
	inboxAddress, _ := models.NewEmailAddress(inbox.EmailAddress)
	emailToSend := models.NewEmail("title", "content")
	receiverAddresses := []models.EmailAddress{*inboxAddress}

	err := emailsImpl.GetEmailClient(publicmocks.EmptyLogger).SendEmails(*emailToSend, receiverAddresses)

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
