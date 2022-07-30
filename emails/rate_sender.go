package emails

import (
	"fmt"
	"gses2.app/api/data"
	"gses2.app/api/exchange_rate"
	"gses2.app/api/smtp"
	"os"
)

// sendEmails sends an email containing a message to each of email addresses.
func sendEmails(addresses []string, message []byte) error {
	// If there's no one to send email to, just return.
	// This counts as a success.
	if len(addresses) == 0 {
		return nil
	}

	// Sender data.
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, addresses, message)
}

// SendRate sends a current exchange rate to a list of email addresses saved in a file storage.
func SendRate() error {
	rate, err := exchange_rate.GetBtcUahRate()
	for err != nil || rate == -1 {
		rate, err = exchange_rate.GetBtcUahRate()
	}

	emailAddresses := data.GetEmailAddresses()
	message := []byte(fmt.Sprintf("Зараз 1 біткоїн коштує %v грн\n", rate)) // convert to string and then to []byte

	return sendEmails(emailAddresses, message)
}
