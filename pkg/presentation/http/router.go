package http

import (
	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/services"
)

func SetupRouter(
	rateService services.ExchangeRateService,
	rateSubscriptionService application.RateSubscriptionService,
	sendRateEmailsService application.SendRateEmailsService,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	handlers := initHandlers(rateService, rateSubscriptionService, sendRateEmailsService)

	registerHandlers(router, handlers)

	return router
}

func initHandlers(
	rateService services.ExchangeRateService,
	rateSubscriptionService application.RateSubscriptionService,
	sendRateEmailsService application.SendRateEmailsService,
) []RequestHandler {
	rateHandler := RateRequestHandler{ExchangeRateService: rateService}
	subscribeHandler := SubscribeRequestHandler{RateSubscriptionService: rateSubscriptionService}
	sendEmailsHandler := SendEmailsRequestHandler{SendRateEmailsService: sendRateEmailsService}

	return []RequestHandler{rateHandler, subscribeHandler, sendEmailsHandler}
}

func registerHandlers(router *gin.Engine, handlers []RequestHandler) {
	for _, handler := range handlers {
		router.Handle(handler.GetMethod(), handler.GetPath(), handler.HandleRequest)
	}
}
