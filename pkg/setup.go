package pkg

import (
	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/domain/services"
	"gses2.app/api/pkg/infrastructure/email"
	"gses2.app/api/pkg/infrastructure/logger"
	"gses2.app/api/pkg/infrastructure/rates"
	"gses2.app/api/pkg/infrastructure/repos"
)

func InitServices() (
	application.BtcToUahRateService,
	application.AddEmailAddressService,
	application.SendBtcToUahRateEmailsService,
) {
	genericExchangeRateService := GetGenericExchangeRateService()
	emailAddressesRepository := repos.GetEmailAddressesFileRepository()
	emailSender := email.GetEmailClient()

	btcToUahService := application.NewBtcToUahServiceImpl(genericExchangeRateService)
	addEmailAddressService := application.NewAddEmailAddressServiceImpl(emailAddressesRepository)
	sendBtcToUahRateEmailsService := application.NewSendBtcToUahRateEmailsServiceImpl(btcToUahService, emailAddressesRepository, emailSender)

	return btcToUahService, addEmailAddressService, sendBtcToUahRateEmailsService
}

func GetGenericExchangeRateService() services.ExchangeRateService {
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
