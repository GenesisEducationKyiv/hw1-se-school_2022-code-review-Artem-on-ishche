package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/services"
)

type SendEmailsRequestHandler struct {
	SendBtcToUahRateEmailsService application.SendBtcToUahRateEmailsService
}

func (handler SendEmailsRequestHandler) GetPath() string {
	return "/sendEmails"
}

func (handler SendEmailsRequestHandler) GetMethod() string {
	return "POST"
}

func (handler SendEmailsRequestHandler) HandleRequest(c *gin.Context) {
	err := handler.SendBtcToUahRateEmailsService.SendBtcToUahRateEmails()
	if errors.Is(err, nil) {
		c.JSON(http.StatusOK, "Success")
	} else if errors.Is(err, services.ErrAPIRequestUnsuccessful) {
		c.JSON(http.StatusBadGateway, "API request has not been successful")
	} else {
		c.JSON(http.StatusInternalServerError, "Some error has occurred")
	}
}
