package routes

import (
	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/presentation/http/handlers"
)

type SubscribeRoute struct {
	handler handlers.SubscribeRequestHandler
}

func (route *SubscribeRoute) GetPath() string {
	return "/subscribe"
}

func (route *SubscribeRoute) GetMethod() string {
	return "POST"
}

func (route *SubscribeRoute) ProcessRequest(ctx *gin.Context) {
	sendJSONResponse(ctx, route.handler.HandleRequest)
}
