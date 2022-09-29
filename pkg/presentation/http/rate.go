package http

import (
	"errors"
	"gses2.app/api/pkg/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/domain/services"
)

type rateRequestParameters struct {
	Base  string `form:"base"`
	Quote string `form:"quote"`
}

type RateRequestHandler struct {
	ExchangeRateService services.ExchangeRateService
}

func (handler RateRequestHandler) GetPath() string {
	return "/rate"
}

func (handler RateRequestHandler) GetMethod() string {
	return "GET"
}

func (handler RateRequestHandler) HandleRequest(c *gin.Context) {
	var params rateRequestParameters

	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, "Input parameters are wrong")

		return
	}

	pair := getCurrencyPair(params.Base, params.Quote)

	exchangeRate, err := handler.ExchangeRateService.GetExchangeRate(*pair)
	if errors.Is(err, nil) {
		c.JSON(http.StatusOK, exchangeRate.Price)
	} else if errors.Is(err, services.ErrAPIRequestUnsuccessful) {
		c.JSON(http.StatusBadGateway, "API request has not been successful")
	} else if errors.Is(err, services.ErrAPIResponseUnmarshallError) {
		c.JSON(http.StatusBadGateway, "API returned unexpected response")
	} else {
		c.JSON(http.StatusInternalServerError, "Some unexpected error has occurred")
	}
}

func getCurrencyPair(baseParam, quoteParam string) *models.CurrencyPair {
	if baseParam == "" {
		baseParam = "btc"
	}

	if quoteParam == "" {
		quoteParam = "uah"
	}

	base := models.NewCurrency(baseParam)
	quote := models.NewCurrency(quoteParam)
	pair := models.NewCurrencyPair(base, quote)

	return &pair
}
