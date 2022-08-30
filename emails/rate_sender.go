package emails

import (
	"fmt"

	"gses2.app/api/config"
	"gses2.app/api/data"
	"gses2.app/api/emails/smtp"
)

func SendRate(rate float64) error {
	emailAddresses := data.GetEmailAddresses()
	message := []byte(fmt.Sprintf("Зараз 1 біткоїн коштує %v грн\n", rate)) // convert to string and then to []byte

	return sendEmails(emailAddresses, message)
}

func sendEmails(addresses []string, message []byte) error {
	if len(addresses) == 0 {
		return nil
	}

	auth := smtp.PlainAuth("", config.EmailAddress, config.EmailPassword, config.SMTPHost)

	return smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.EmailAddress, addresses, message)
}
