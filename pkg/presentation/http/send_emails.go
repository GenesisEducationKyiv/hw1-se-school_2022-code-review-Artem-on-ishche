package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/application"
)

type sendEmailsRequestParameters struct {
	Key   string `form:"key" binding:"required"`
	Base  string `form:"base"`
	Quote string `form:"quote"`
}

type SendEmailsRequestHandler struct {
	SendRateEmailsService application.SendRateEmailsService
}

func (handler SendEmailsRequestHandler) GetPath() string {
	return "/sendEmails"
}

func (handler SendEmailsRequestHandler) GetMethod() string {
	return "POST"
}

func (handler SendEmailsRequestHandler) HandleRequest(ctx *gin.Context) {
	var params sendEmailsRequestParameters

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, "Input parameters are wrong")

		return
	}

	key := params.Key
	pair := getCurrencyPair(params.Base, params.Quote)

	err := handler.SendRateEmailsService.SendRateEmails(pair, key)
	if errors.Is(err, nil) {
		ctx.JSON(http.StatusOK, "Success")
	} else if errors.Is(err, application.ErrValidationError) {
		ctx.JSON(http.StatusUnauthorized, "Provided key is not valid")
	} else if errors.Is(err, application.ErrAPIRequestUnsuccessful) {
		ctx.JSON(http.StatusBadGateway, "API request has not been successful")
	} else {
		ctx.JSON(http.StatusInternalServerError, "Some error has occurred")
	}
}
