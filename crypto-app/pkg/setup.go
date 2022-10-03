package pkg

import (
	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/domain/services"
	"gses2.app/api/pkg/infrastructure/email"
	"gses2.app/api/pkg/infrastructure/rates"
	"gses2.app/api/pkg/infrastructure/repos"
)

func InitServices(loggerService services.Logger) (
	services.ExchangeRateService,
	application.RateSubscriptionService,
	application.SendRateEmailsService,
	services.Logger,
) {

	genericExchangeRateService := GetGenericExchangeRateService(loggerService)
	repositoryGetter := repos.NewEmailAddressesFileRepoGetter(loggerService)
	emailSender := email.GetEmailClient(loggerService)

	subscribeToRateService := application.NewSubscribeToRateServiceImpl(repositoryGetter, loggerService)
	sendBtcToUahRateEmailsService := application.NewSendRateEmailsServiceImpl(
		config.AdminKey,
		genericExchangeRateService,
		repositoryGetter,
		emailSender,
		loggerService,
	)

	return genericExchangeRateService, subscribeToRateService, sendBtcToUahRateEmailsService, loggerService
}

func GetGenericExchangeRateService(loggerService services.Logger) services.ExchangeRateService {
	fiveMinutes := 5.0
	cacherRateService := rates.CacherRateServiceFactory{MaxTime: fiveMinutes, Logger: loggerService}.CreateRateService()

	coinRateService := rates.CoinAPIClientFactory{Cacher: cacherRateService, Logger: loggerService}.CreateRateService()
	nomicsRateService := rates.NomicsAPIClientFactory{Cacher: cacherRateService, Logger: loggerService}.CreateRateService()
	binanceRateService := rates.BinanceAPIClientFactory{Cacher: cacherRateService, Logger: loggerService}.CreateRateService()

	switch config.CryptoCurrencyProvider {
	case "coin":
		cacherRateService.SetNext(&coinRateService)
		coinRateService.SetNext(&nomicsRateService)
		nomicsRateService.SetNext(&binanceRateService)

		return cacherRateService
	case "nomics":
		cacherRateService.SetNext(&nomicsRateService)
		nomicsRateService.SetNext(&coinRateService)
		coinRateService.SetNext(&binanceRateService)

		return cacherRateService
	case "binance":
		cacherRateService.SetNext(&binanceRateService)
		binanceRateService.SetNext(&coinRateService)
		coinRateService.SetNext(&nomicsRateService)

		return cacherRateService
	default:
		panic("Wrong crypto provider .env value")
	}
}
