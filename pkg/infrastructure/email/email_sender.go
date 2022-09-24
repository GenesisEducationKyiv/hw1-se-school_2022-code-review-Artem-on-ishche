package email

import (
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
}

func GetEmailClient() services.EmailSender {
	return &emailClient{
		emailAddress:  config.EmailAddress,
		emailPassword: config.EmailPassword,
		smtpHost:      "smtp.gmail.com",
		smtpPort:      "587",
	}
}

func (emailClient *emailClient) SendEmails(email models.Email, receiverAddresses []string) error {
	if len(receiverAddresses) == 0 {
		return nil
	}

	auth := smtp.PlainAuth("", emailClient.emailAddress, emailClient.emailPassword, emailClient.smtpHost)
	message := []byte(email.Body)

	return smtp.SendMail(emailClient.smtpHost+":"+emailClient.smtpPort, auth, emailClient.emailAddress, receiverAddresses, message)
}
