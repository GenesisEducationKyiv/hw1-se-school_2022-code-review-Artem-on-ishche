package http

import (
	"gses2.app/api/pkg/domain/models"
	httpPresentation "gses2.app/api/pkg/presentation/http/routes"
)

type addEmailAddressFunction func(emailAddress models.EmailAddress) error

var addEmailAddressTestFunction addEmailAddressFunction

type rateSubscriptionServiceTestDouble struct{}

func (service rateSubscriptionServiceTestDouble) Subscribe(emailAddress *models.EmailAddress, _ *models.CurrencyPair) error {
	return addEmailAddressTestFunction(*emailAddress)
}

var testSubscribeRequestHandler = httpPresentation.SubscribeRequestHandler{RateSubscriptionService: rateSubscriptionServiceTestDouble{}}
