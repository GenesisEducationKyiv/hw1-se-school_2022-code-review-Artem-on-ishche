package services

import "fmt"

func SendBtcToUahRateEmails(rateService ExchangeRateService, storage EmailAddressesStorage, sender EmailSender) error {
	rate, err := GetBtcToUahRate(rateService)
	if err != nil {
		return err
	}

	email := getEmailWithRate(rate)
	receiverAddresses := storage.GetEmailAddresses()

	return sender.SendEmails(email, receiverAddresses)
}

func getEmailWithRate(rate float64) Email {
	title := "BTC to UAH rate"
	body := fmt.Sprintf("Зараз 1 біткоїн коштує %v грн\n", rate)

	return *NewEmail(title, body)
}
