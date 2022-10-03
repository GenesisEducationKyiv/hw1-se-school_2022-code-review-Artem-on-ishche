package routes

import (
	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/services"
	"gses2.app/api/pkg/presentation/http/handlers"
)

func SetupRouter(
	rateService services.ExchangeRateService,
	rateSubscriptionService application.RateSubscriptionService,
	sendRateEmailsService application.SendRateEmailsService,
	logger services.Logger,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	routes := initRoutes(rateService, rateSubscriptionService, sendRateEmailsService, logger)

	registerRoutes(router, routes)

	return router
}

func initRoutes(
	rateService services.ExchangeRateService,
	rateSubscriptionService application.RateSubscriptionService,
	sendRateEmailsService application.SendRateEmailsService,
	logger services.Logger,
) []RequestRoute {
	rateRoute := RateRoute{
		handler: handlers.RateRequestHandler{ExchangeRateService: rateService, Logger: logger},
		logger:  logger,
	}
	subscribeRoute := SubscribeRoute{
		handler: handlers.SubscribeRequestHandler{RateSubscriptionService: rateSubscriptionService, Logger: logger},
		logger:  logger,
	}
	sendEmailsRoute := SendEmailsRoute{
		handler: handlers.SendEmailsRequestHandler{SendRateEmailsService: sendRateEmailsService, Logger: logger},
		logger:  logger,
	}

	return []RequestRoute{&rateRoute, &subscribeRoute, &sendEmailsRoute}
}

func registerRoutes(router *gin.Engine, routes []RequestRoute) {
	for _, route := range routes {
		router.Handle(route.GetMethod(), route.GetPath(), route.ProcessRequest)
	}
}
