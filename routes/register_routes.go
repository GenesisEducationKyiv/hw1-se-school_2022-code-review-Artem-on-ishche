package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gses2.app/api/config"
	"gses2.app/api/handlers"
	"gses2.app/api/implementations"
)

func StartRouter() {
	router := mux.NewRouter().StrictSlash(true)

	registerRoutes(router)

	log.Fatal(http.ListenAndServe(config.NetworkPort, router))
}

func registerRoutes(router *mux.Router) {
	router.HandleFunc("/rate", rateRoute).Methods("GET")
	router.HandleFunc("/subscribe", subscribeRoute).Methods("POST")
	router.HandleFunc("/sendEmails", sendEmailsRoute).Methods("POST")
}

func rateRoute(responseWriter http.ResponseWriter, request *http.Request) {
	genericExchangeRateService := implementations.GetExchangeRateCoinAPIClient()
	rateRequestHandler := handlers.NewBtcToUahRateRequestHandler(genericExchangeRateService)

	handleRoute(responseWriter, request, rateRequestHandler)
}

func subscribeRoute(responseWriter http.ResponseWriter, request *http.Request) {
	emailAddressesStorage := implementations.GetEmailAddressesFileStorage()
	subscribeRequestHandler := handlers.NewSubscribeRequestHandler(emailAddressesStorage)

	handleRoute(responseWriter, request, subscribeRequestHandler)
}

func sendEmailsRoute(responseWriter http.ResponseWriter, request *http.Request) {
	genericExchangeRateService := implementations.GetExchangeRateCoinAPIClient()
	emailAddressesStorage := implementations.GetEmailAddressesFileStorage()
	emailSender := implementations.GetEmailClient()
	sendEmailsRequestHandler := handlers.NewSendEmailsRequestHandler(genericExchangeRateService, emailAddressesStorage, emailSender)

	handleRoute(responseWriter, request, sendEmailsRequestHandler)
}

func handleRoute(responseWriter http.ResponseWriter, request *http.Request, handler handlers.RequestHandler) {
	response := handler.HandleRequest(request)
	responseSender.sendResponse(responseWriter, response.StatusCode, response.Message)
}
