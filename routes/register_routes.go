package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gses2.app/api/config"
	"gses2.app/api/handlers"
	"gses2.app/api/implementations/email"
	"gses2.app/api/implementations/logger"
	"gses2.app/api/implementations/rates"
	"gses2.app/api/implementations/repos"
	"gses2.app/api/services"
)

func StartRouter() {
	router := mux.NewRouter().StrictSlash(true)
	routes := initRoutes()

	registerRoutes(router, routes)

	log.Fatal(http.ListenAndServe(config.NetworkPort, router))
}

func initRoutes() []route {
	genericExchangeRateService := getGenericExchangeRateService()
	emailAddressesRepository := repos.GetEmailAddressesFileRepository()
	emailSender := email.GetEmailClient()

	rateRoute := initRateRoute(&genericExchangeRateService)
	subscribeRoute := initSubscribeRoute(&emailAddressesRepository)
	sendEmailsRoute := initSendEmailsRoute(&genericExchangeRateService, &emailAddressesRepository, &emailSender)

	return []route{rateRoute, subscribeRoute, sendEmailsRoute}
}

func initRateRoute(genericExchangeRateService *services.ExchangeRateService) route {
	rateRequestHandler := handlers.NewBtcToUahRateRequestHandler(*genericExchangeRateService)

	return route{
		path:    "/rate",
		method:  "GET",
		handler: &rateRequestHandler,
	}
}

func initSubscribeRoute(emailAddressesRepository *services.EmailAddressesRepository) route {
	subscribeRequestHandler := handlers.NewSubscribeRequestHandler(*emailAddressesRepository)

	return route{
		path:    "/subscribe",
		method:  "POST",
		handler: &subscribeRequestHandler,
	}
}

func initSendEmailsRoute(
	genericExchangeRateService *services.ExchangeRateService,
	emailAddressesRepository *services.EmailAddressesRepository,
	emailSender *services.EmailSender,
) route {
	sendEmailsRequestHandler := handlers.NewSendEmailsRequestHandler(*genericExchangeRateService, *emailAddressesRepository, *emailSender)

	return route{
		path:    "/sendEmails",
		method:  "POST",
		handler: &sendEmailsRequestHandler,
	}
}

func getGenericExchangeRateService() services.ExchangeRateService {
	fiveMinutes := 5.0
	cacherRateService := rates.CacherRateServiceFactory{MaxTime: fiveMinutes}.CreateRateService()
	loggerService := logger.ConsoleLogger{}

	mediator := getMediator(cacherRateService, loggerService)

	coinRateService := rates.CoinAPIClientFactory{Mediator: mediator}.CreateRateService()
	nomicsRateService := rates.NomicsAPIClientFactory{Mediator: mediator}.CreateRateService()
	coinMarketCapRateService := rates.CoinMarketCapAPIClientFactory{Mediator: mediator}.CreateRateService()

	switch config.CryptoCurrencyProvider {
	case "coin":
		cacherRateService.SetNext(&coinRateService)
		coinRateService.SetNext(&nomicsRateService)
		nomicsRateService.SetNext(&coinMarketCapRateService)

		return cacherRateService
	case "nomics":
		cacherRateService.SetNext(&nomicsRateService)
		nomicsRateService.SetNext(&coinRateService)
		coinRateService.SetNext(&coinMarketCapRateService)

		return cacherRateService
	case "coin_market_cap":
		cacherRateService.SetNext(&coinMarketCapRateService)
		coinMarketCapRateService.SetNext(&coinRateService)
		coinRateService.SetNext(&nomicsRateService)

		return cacherRateService
	default:
		panic("Wrong crypto provider .env value")
	}
}

func getMediator(cacherRateService rates.CacherRateService, loggerService services.Logger) *rates.Mediator {
	mediator := rates.NewMediator()

	err := mediator.Attach(rates.NewRateReturnedObserver{Cacher: cacherRateService}, rates.NewRateReturnedEvent{}.GetName())
	if err != nil {
		return nil
	}

	err = mediator.Attach(rates.FailureAPIResponseReceivedObserver{Logger: loggerService}, rates.FailureAPIResponseReceivedEvent{}.GetName())
	if err != nil {
		return nil
	}

	err = mediator.Attach(rates.SuccessAPIResponseReceivedObserver{Logger: loggerService}, rates.SuccessAPIResponseReceivedEvent{}.GetName())
	if err != nil {
		return nil
	}

	return &mediator
}

func registerRoutes(router *mux.Router, routes []route) {
	for _, route := range routes {
		router.HandleFunc(route.path, route.processRequest).Methods(route.method)
	}
}
