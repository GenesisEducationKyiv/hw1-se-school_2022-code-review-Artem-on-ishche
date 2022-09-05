package implementations

import (
	"gses2.app/api/config"
	"gses2.app/api/implementations/smtp"
	"gses2.app/api/services"
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
		smtpHost:      config.SMTPHost,
		smtpPort:      config.SMTPPort,
	}
}

func (emailClient *emailClient) SendEmails(email services.Email, receiverAddresses []string) error {
	if len(receiverAddresses) == 0 {
		return nil
	}

	auth := smtp.PlainAuth("", emailClient.emailAddress, emailClient.emailPassword, emailClient.smtpHost)
	message := []byte(email.Body)

	return smtp.SendMail(emailClient.smtpHost+":"+emailClient.smtpPort, auth, emailClient.emailAddress, receiverAddresses, message)
}
