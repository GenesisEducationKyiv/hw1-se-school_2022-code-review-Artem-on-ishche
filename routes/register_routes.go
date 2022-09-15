package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gses2.app/api/config"
	"gses2.app/api/handlers"
	"gses2.app/api/implementations/email"
	"gses2.app/api/implementations/rates"
	"gses2.app/api/implementations/repos"
	"gses2.app/api/services"
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
	genericExchangeRateService := getGenericExchangeRateService()
	rateRequestHandler := handlers.NewBtcToUahRateRequestHandler(genericExchangeRateService)

	handleRoute(responseWriter, request, rateRequestHandler)
}

func subscribeRoute(responseWriter http.ResponseWriter, request *http.Request) {
	emailAddressesStorage := repos.GetEmailAddressesFileRepository()
	subscribeRequestHandler := handlers.NewSubscribeRequestHandler(emailAddressesStorage)

	handleRoute(responseWriter, request, subscribeRequestHandler)
}

func sendEmailsRoute(responseWriter http.ResponseWriter, request *http.Request) {
	genericExchangeRateService := getGenericExchangeRateService()
	emailAddressesStorage := repos.GetEmailAddressesFileRepository()
	emailSender := email.GetEmailClient()
	sendEmailsRequestHandler := handlers.NewSendEmailsRequestHandler(genericExchangeRateService, emailAddressesStorage, emailSender)

	handleRoute(responseWriter, request, sendEmailsRequestHandler)
}

func handleRoute(responseWriter http.ResponseWriter, request *http.Request, handler handlers.RequestHandler) {
	response := handler.HandleRequest(request)
	responseSender.sendResponse(responseWriter, response.StatusCode, response.Message)
}

func getGenericExchangeRateService() services.ExchangeRateService {
	coinRateService := rates.CoinAPIClientFactory{}.CreateRateService()
	nomicsRateService := rates.NomicsAPIClientFactory{}.CreateRateService()
	coinMarketCapRateService := rates.CoinMarketCapAPIClientFactory{}.CreateRateService()

	coinRateService.SetNext(&nomicsRateService)
	nomicsRateService.SetNext(&coinMarketCapRateService)

	switch config.CryptoCurrencyProvider {
	case "coin":
		return coinRateService
	case "nomics":
		return nomicsRateService
	case "coin_market_cap":
		return coinMarketCapRateService
	default:
		panic("Wrong crypto provider .env value")
	}
}
