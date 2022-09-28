package http

import httpPresentation "gses2.app/api/pkg/presentation/http"

type sendEmailsFunction func() error

var sendBtcToUahRateEmailsTestFunction sendEmailsFunction

type sendBtcToUahRateEmailsServiceTestDouble struct{}

func (service sendBtcToUahRateEmailsServiceTestDouble) SendBtcToUahRateEmails() error {
	return sendBtcToUahRateEmailsTestFunction()
}

var testSendEmailsHandler = httpPresentation.SendEmailsRequestHandler{
	SendBtcToUahRateEmailsService: sendBtcToUahRateEmailsServiceTestDouble{},
}
