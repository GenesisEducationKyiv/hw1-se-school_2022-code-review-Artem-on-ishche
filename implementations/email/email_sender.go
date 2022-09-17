package email

import (
	"gses2.app/api/config"
	smtp2 "gses2.app/api/implementations/email/smtp"
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
		smtpHost:      "smtp.gmail.com",
		smtpPort:      "587",
	}
}

func (emailClient *emailClient) SendEmails(email services.Email, receiverAddresses []string) error {
	if len(receiverAddresses) == 0 {
		return nil
	}

	auth := smtp2.PlainAuth("", emailClient.emailAddress, emailClient.emailPassword, emailClient.smtpHost)
	message := []byte(email.Body)

	return smtp2.SendMail(emailClient.smtpHost+":"+emailClient.smtpPort, auth, emailClient.emailAddress, receiverAddresses, message)
}
