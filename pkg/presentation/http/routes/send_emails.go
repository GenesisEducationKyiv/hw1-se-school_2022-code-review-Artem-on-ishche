package routes

import (
	"github.com/gin-gonic/gin"
	"gses2.app/api/pkg/presentation/http/handlers"
)

type SendEmailsRoute struct {
	handler handlers.SendEmailsRequestHandler
}

func (route *SendEmailsRoute) GetPath() string {
	return "/sendEmails"
}

func (route *SendEmailsRoute) GetMethod() string {
	return "POST"
}

func (route *SendEmailsRoute) ProcessRequest(ctx *gin.Context) {
	sendJSONResponse(ctx, route.handler.HandleRequest)
}
