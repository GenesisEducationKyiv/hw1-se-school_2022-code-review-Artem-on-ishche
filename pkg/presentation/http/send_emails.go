package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/services"
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

func (handler SendEmailsRequestHandler) HandleRequest(c *gin.Context) {
	var params sendEmailsRequestParameters

	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, "Input parameters are wrong")

		return
	}

	key := params.Key
	pair := getCurrencyPair(params.Base, params.Quote)

	err := handler.SendRateEmailsService.SendRateEmails(pair, key)
	if errors.Is(err, nil) {
		c.JSON(http.StatusOK, "Success")
	} else if errors.Is(err, application.ErrValidationError) {
		c.JSON(http.StatusUnauthorized, "Provided key is not valid")
	} else if errors.Is(err, services.ErrAPIRequestUnsuccessful) {
		c.JSON(http.StatusBadGateway, "API request has not been successful")
	} else {
		c.JSON(http.StatusInternalServerError, "Some error has occurred")
	}
}
