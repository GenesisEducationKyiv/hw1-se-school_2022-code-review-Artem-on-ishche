package email

import (
	"fmt"

	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
	"gses2.app/api/pkg/infrastructure/email/smtp"
)

type emailClient struct {
	emailAddress  string
	emailPassword string
	smtpHost      string
	smtpPort      string

	logger services.Logger
}

func GetEmailClient(logger services.Logger) services.EmailSender {
	return &emailClient{
		emailAddress:  config.EmailAddress,
		emailPassword: config.EmailPassword,
		smtpHost:      "smtp.gmail.com",
		smtpPort:      "587",

		logger: logger,
	}
}

func (emailClient *emailClient) SendEmails(email models.EmailMessage, receiverAddresses []models.EmailAddress) error {
	emailClient.logger.Debug(fmt.Sprintf("SendEmails() called with email=%v, receiverAddresses=%v", email, receiverAddresses))

	if len(receiverAddresses) == 0 {
		emailClient.logger.Debug("No one to send emails to, success")

		return nil
	}

	auth := smtp.PlainAuth("", emailClient.emailAddress, emailClient.emailPassword, emailClient.smtpHost)
	receiverAddressStrings := getEmailAddressStrings(receiverAddresses)
	message := []byte(email.Body)

	return smtp.SendMail(emailClient.smtpHost+":"+emailClient.smtpPort, auth, emailClient.emailAddress, receiverAddressStrings, message)
}

func getEmailAddressStrings(addresses []models.EmailAddress) []string {
	addressStrings := make([]string, len(addresses))

	for i, addr := range addresses {
		addressStrings[i] = addr.String()
	}

	return addressStrings
}
