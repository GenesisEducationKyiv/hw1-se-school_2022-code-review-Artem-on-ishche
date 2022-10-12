package routes

import (
	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/domain/services"
	"gses2.app/api/pkg/presentation/http/handlers"
)

type SendEmailsRoute struct {
	handler handlers.SendEmailsRequestHandler
	logger  services.Logger
}

func (route *SendEmailsRoute) GetPath() string {
	return "/sendEmails"
}

func (route *SendEmailsRoute) GetMethod() string {
	return "POST"
}

func (route *SendEmailsRoute) ProcessRequest(ctx *gin.Context) {
	route.logger.Info(route.GetPath() + "?" + ctx.Request.URL.RawQuery + " called")

	sendJSONResponse(route.logger, ctx, route.handler.HandleRequest)
}
