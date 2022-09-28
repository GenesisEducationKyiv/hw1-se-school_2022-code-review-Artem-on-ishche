package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/services"
)

type BtcToUahRateRequestHandler struct {
	BtcToUahService application.BtcToUahRateService
}

func (handler BtcToUahRateRequestHandler) GetPath() string {
	return "/rate"
}

func (handler BtcToUahRateRequestHandler) GetMethod() string {
	return "GET"
}

func (handler BtcToUahRateRequestHandler) HandleRequest(c *gin.Context) {
	exchangeRate, err := handler.BtcToUahService.GetBtcToUahRate()
	if errors.Is(err, nil) {
		c.JSON(http.StatusOK, exchangeRate.Price)
	} else if errors.Is(err, services.ErrAPIRequestUnsuccessful) {
		c.JSON(http.StatusBadGateway, "API request has not been successful")
	} else {
		c.JSON(http.StatusInternalServerError, "Some unexpected error has occurred")
	}
}
