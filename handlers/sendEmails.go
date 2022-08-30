package handlers

import (
	"net/http"

	"gses2.app/api/emails"
	"gses2.app/api/rate"
)

func SendEmailsHandler(responseWriter http.ResponseWriter, _ *http.Request) {
	exchangeRate, err := rate.GetBtcToUahRate()
	if isRateWrong(exchangeRate, err) {
		sendBadRequestResponse(responseWriter, "An error has occurred")

		return
	}

	err = emails.SendRate(exchangeRate)
	if err != nil {
		sendInternalServerErrorResponse(responseWriter)

		return
	}

	sendSuccessResponse(responseWriter, "Success")
}
