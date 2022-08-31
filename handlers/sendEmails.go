package handlers

import (
	"net/http"

	"gses2.app/api/emails"
	"gses2.app/api/rates"
)

func sendEmailsHandler(responseWriter http.ResponseWriter, _ *http.Request) {
	exchangeRate, err := rates.GetBtcToUahRate()
	if isRateWrong(exchangeRate, err) {
		sendBadRequestResponse(responseWriter, "An error has occurred")

		return
	}

	err = emails.SendEmailsAboutRate(exchangeRate)
	if err != nil {
		sendInternalServerErrorResponse(responseWriter)

		return
	}

	sendSuccessResponse(responseWriter, "Success")
}
