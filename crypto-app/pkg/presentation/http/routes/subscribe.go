package routes

import (
	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/domain/services"
	"gses2.app/api/pkg/presentation/http/handlers"
)

type SubscribeRoute struct {
	handler handlers.SubscribeRequestHandler
	logger  services.Logger
}

func (route *SubscribeRoute) GetPath() string {
	return "/subscribe"
}

func (route *SubscribeRoute) GetMethod() string {
	return "POST"
}

func (route *SubscribeRoute) ProcessRequest(ctx *gin.Context) {
	route.logger.Info(ctx.Request.URL.RawQuery + " called")

	sendJSONResponse(route.logger, ctx, route.handler.HandleRequest)
}
