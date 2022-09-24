package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/config"
)

func StartRouter(
	btcToUahService application.BtcToUahRateService,
	addEmailAddressService application.AddEmailAddressService,
	sendBtcToUahRateEmailsService application.SendBtcToUahRateEmailsService,
) {
	router := mux.NewRouter().StrictSlash(true)
	handlers := InitHandlers(btcToUahService, addEmailAddressService, sendBtcToUahRateEmailsService)

	registerHandlers(router, handlers)

	log.Fatal(http.ListenAndServe(config.NetworkPort, router))
}

func InitHandlers(
	btcToUahService application.BtcToUahRateService,
	addEmailAddressService application.AddEmailAddressService,
	sendBtcToUahRateEmailsService application.SendBtcToUahRateEmailsService,
) []RequestHandler {
	rateHandler := BtcToUahRateRequestHandler{BtcToUahService: btcToUahService}
	subscribeHandler := SubscribeRequestHandler{AddEmailAddressService: addEmailAddressService}
	sendEmailsHandler := SendEmailsRequestHandler{SendBtcToUahRateEmailsService: sendBtcToUahRateEmailsService}

	return []RequestHandler{rateHandler, subscribeHandler, sendEmailsHandler}
}

func registerHandlers(router *mux.Router, handlers []RequestHandler) {
	for _, handler := range handlers {
		router.HandleFunc(handler.GetPath(), GetHandlerFunction(handler)).Methods(handler.GetMethod())
	}
}
