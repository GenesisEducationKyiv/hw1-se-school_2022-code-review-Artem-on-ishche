package main

import (
	"fmt"
	"net/smtp"
	"os"
)

// sendEmails sends an email containing a message to each of email addresses.
func sendEmails(addresses []string, message []byte) {
	// Sender data.
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, addresses, message)
	for err != nil {
		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, addresses, message)
	}
}

// sendRate sends a current exchange rate to a list of email addresses saved in a file storage.
func sendRate() {
	rate, err := getBtcUahRate()
	for err != nil || rate == -1 {
		rate, err = getBtcUahRate()
	}

	emailAddresses := getEmailAddresses()
	message := []byte(fmt.Sprintf("Зараз 1 біткоїн коштує %v грн\n", rate)) // convert to string and then to []byte

	sendEmails(emailAddresses, message)
}
