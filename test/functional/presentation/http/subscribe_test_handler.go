package http

import (
	"gses2.app/api/pkg/domain/models"
	httpPresentation "gses2.app/api/pkg/presentation/http"
)

type addEmailAddressFunction func(emailAddress models.EmailAddress) error

var addEmailAddressTestFunction addEmailAddressFunction

type addEmailAddressServiceTestDouble struct{}

func (service addEmailAddressServiceTestDouble) AddEmailAddress(emailAddress models.EmailAddress) error {
	return addEmailAddressTestFunction(emailAddress)
}

var testSubscribeRequestHandler = httpPresentation.SubscribeRequestHandler{AddEmailAddressService: addEmailAddressServiceTestDouble{}}
