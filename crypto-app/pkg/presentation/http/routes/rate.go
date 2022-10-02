package routes

import (
	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/domain/services"
	"gses2.app/api/pkg/presentation/http/handlers"
)

type RateRoute struct {
	handler handlers.RateRequestHandler
	logger  services.Logger
}

func (route *RateRoute) GetPath() string {
	return "/rate"
}

func (route *RateRoute) GetMethod() string {
	return "GET"
}

func (route *RateRoute) ProcessRequest(ctx *gin.Context) {
	route.logger.Info(ctx.Request.URL.RawQuery + " called")

	sendJSONResponse(route.logger, ctx, route.handler.HandleRequest)
}
