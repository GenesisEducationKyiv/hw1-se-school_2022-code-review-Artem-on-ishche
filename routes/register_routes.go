package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gses2.app/api/config"
	"gses2.app/api/emails"
	"gses2.app/api/handlers"
	"gses2.app/api/rates"
	"gses2.app/api/storage"
)

func RegisterRoutes() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/rate", rateRoute).Methods("GET")
	router.HandleFunc("/subscribe", subscribeRoute).Methods("POST")
	router.HandleFunc("/sendEmails", sendEmailsRoute).Methods("POST")

	log.Fatal(http.ListenAndServe(config.NetworkPort, router))
}

func rateRoute(responseWriter http.ResponseWriter, request *http.Request) {
	handlers.ExchangeRateServiceImpl = &rates.ExchangeRateCoinApiClient
	registerRoute(responseWriter, request, handlers.RateRequestHandler)
}

func subscribeRoute(responseWriter http.ResponseWriter, request *http.Request) {
	handlers.EmailAddressesStorageImpl = data.EmailAddressesFileStorage
	registerRoute(responseWriter, request, handlers.SubscribeRequestHandler)
}

func sendEmailsRoute(responseWriter http.ResponseWriter, request *http.Request) {
	handlers.EmailSenderImpl = &emails.EmailClient
	registerRoute(responseWriter, request, handlers.SendEmailsRequestHandler)
}

func registerRoute(responseWriter http.ResponseWriter, request *http.Request, handler handlers.RequestHandler) {
	response := handler.HandleRequest(request)
	responseSender.sendResponse(responseWriter, response.StatusCode, response.Message)
}
