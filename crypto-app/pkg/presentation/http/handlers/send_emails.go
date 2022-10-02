package handlers

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

func (handler *SendEmailsRequestHandler) HandleRequest(ctx *gin.Context) *JSONResponse {
	var params sendEmailsRequestParameters

	if err := ctx.ShouldBind(&params); err != nil {
		return NewJSONResponse(http.StatusBadRequest, "Input parameters are wrong")
	}

	key := params.Key
	pair := getCurrencyPair(params.Base, params.Quote)

	err := handler.SendRateEmailsService.SendRateEmails(pair, key)
	if errors.Is(err, nil) {
		return NewJSONResponse(http.StatusOK, "Success")
	} else if errors.Is(err, application.ErrValidationError) {
		return NewJSONResponse(http.StatusUnauthorized, "Provided key is not valid")
	} else if errors.Is(err, application.ErrAPIRequestUnsuccessful) {
		return NewJSONResponse(http.StatusBadGateway, "API request has not been successful")
	} else {
		return NewJSONResponse(http.StatusInternalServerError, "Some error has occurred")
	}
}
