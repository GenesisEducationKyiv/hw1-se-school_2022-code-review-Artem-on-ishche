package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gses2.app/api/config"
	"gses2.app/api/handlers"
	"gses2.app/api/implementations"
	"gses2.app/api/services"
)

func RegisterRoutes() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/rate", rateRoute).Methods("GET")
	router.HandleFunc("/subscribe", subscribeRoute).Methods("POST")
	router.HandleFunc("/sendEmails", sendEmailsRoute).Methods("POST")

	log.Fatal(http.ListenAndServe(config.NetworkPort, router))
}

func rateRoute(responseWriter http.ResponseWriter, request *http.Request) {
	genericExchangeRateService := implementations.GetExchangeRateCoinApiClient()
	btcToUahRateService := services.NewBtcToUahServiceImpl(genericExchangeRateService)
	rateRequestHandler := handlers.NewBtcToUahRateRequestHandler(btcToUahRateService)

	registerRoute(responseWriter, request, rateRequestHandler)
}

func subscribeRoute(responseWriter http.ResponseWriter, request *http.Request) {
	emailAddressesStorage := implementations.GetEmailAddressesFileStorage()
	addEmailAddressService := services.NewAddEmailAddressServiceImpl(emailAddressesStorage)
	subscribeRequestHandler := handlers.NewSubscribeRequestHandler(addEmailAddressService)

	registerRoute(responseWriter, request, subscribeRequestHandler)
}

func sendEmailsRoute(responseWriter http.ResponseWriter, request *http.Request) {
	genericExchangeRateService := implementations.GetExchangeRateCoinApiClient()
	btcToUahRateService := services.NewBtcToUahServiceImpl(genericExchangeRateService)
	emailAddressesStorage := implementations.GetEmailAddressesFileStorage()
	emailSender := implementations.GetEmailClient()
	sendBtcToUahRateEmailsService := services.NewSendBtcToUahRateEmailsServiceImpl(btcToUahRateService, emailAddressesStorage, emailSender)
	sendEmailsRequestHandler := handlers.NewSendEmailsRequestHandler(sendBtcToUahRateEmailsService)

	registerRoute(responseWriter, request, sendEmailsRequestHandler)
}

func registerRoute(responseWriter http.ResponseWriter, request *http.Request, handler handlers.RequestHandler) {
	response := handler.HandleRequest(request)
	responseSender.sendResponse(responseWriter, response.StatusCode, response.Message)
}
