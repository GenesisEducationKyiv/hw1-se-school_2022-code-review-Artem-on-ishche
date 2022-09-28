package http

import (
	"github.com/gin-gonic/gin"
	"gses2.app/api/pkg/application"
)

func SetupRouter(
	btcToUahService application.BtcToUahRateService,
	addEmailAddressService application.AddEmailAddressService,
	sendBtcToUahRateEmailsService application.SendBtcToUahRateEmailsService,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	handlers := initHandlers(btcToUahService, addEmailAddressService, sendBtcToUahRateEmailsService)

	registerHandlers(router, handlers)

	return router
}

func initHandlers(
	btcToUahService application.BtcToUahRateService,
	addEmailAddressService application.AddEmailAddressService,
	sendBtcToUahRateEmailsService application.SendBtcToUahRateEmailsService,
) []RequestHandler {
	rateHandler := BtcToUahRateRequestHandler{BtcToUahService: btcToUahService}
	subscribeHandler := SubscribeRequestHandler{AddEmailAddressService: addEmailAddressService}
	sendEmailsHandler := SendEmailsRequestHandler{SendBtcToUahRateEmailsService: sendBtcToUahRateEmailsService}

	return []RequestHandler{rateHandler, subscribeHandler, sendEmailsHandler}
}

func registerHandlers(router *gin.Engine, handlers []RequestHandler) {
	for _, handler := range handlers {
		router.Handle(handler.GetMethod(), handler.GetPath(), handler.HandleRequest)
	}
}
