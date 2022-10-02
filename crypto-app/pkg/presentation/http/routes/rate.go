package routes

import (
	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/presentation/http/handlers"
)

type RateRoute struct {
	handler handlers.RateRequestHandler
}

func (route *RateRoute) GetPath() string {
	return "/rate"
}

func (route *RateRoute) GetMethod() string {
	return "GET"
}

func (route *RateRoute) ProcessRequest(ctx *gin.Context) {
	sendJSONResponse(ctx, route.handler.HandleRequest)
}
