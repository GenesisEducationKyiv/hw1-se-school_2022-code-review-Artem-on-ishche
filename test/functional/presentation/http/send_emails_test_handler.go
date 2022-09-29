package http

import (
	"gses2.app/api/pkg/domain/models"
	httpPresentation "gses2.app/api/pkg/presentation/http"
)

type sendEmailsFunction func() error

var sendBtcToUahRateEmailsTestFunction sendEmailsFunction

type sendRateEmailsServiceTestDouble struct{}

func (service sendRateEmailsServiceTestDouble) SendRateEmails(*models.CurrencyPair, string) error {
	return sendBtcToUahRateEmailsTestFunction()
}

var testSendEmailsHandler = httpPresentation.SendEmailsRequestHandler{
	SendRateEmailsService: sendRateEmailsServiceTestDouble{},
}
