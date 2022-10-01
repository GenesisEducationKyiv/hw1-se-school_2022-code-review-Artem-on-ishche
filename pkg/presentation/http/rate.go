package http

import (
	"errors"
	"gses2.app/api/pkg/application"
	"net/http"

	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

const (
	defaultBaseName  = "BTC"
	defaultQuoteName = "UAH"
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

func (handler RateRequestHandler) HandleRequest(ctx *gin.Context) {
	var params rateRequestParameters

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, "Input parameters are wrong")

		return
	}

	handleEmptyParameter(&params.Base, defaultBaseName)
	handleEmptyParameter(&params.Quote, defaultQuoteName)

	pair := getCurrencyPair(params.Base, params.Quote)

	exchangeRate, err := handler.ExchangeRateService.GetExchangeRate(*pair)
	if errors.Is(err, nil) {
		ctx.JSON(http.StatusOK, exchangeRate.Price)
	} else if errors.Is(err, application.ErrAPIRequestUnsuccessful) {
		ctx.JSON(http.StatusBadGateway, "API request has not been successful")
	} else if errors.Is(err, application.ErrAPIResponseUnmarshallError) {
		ctx.JSON(http.StatusBadGateway, "API returned unexpected response")
	} else {
		ctx.JSON(http.StatusInternalServerError, "Some unexpected error has occurred")
	}
}

func handleEmptyParameter(param *string, defaultValue string) {
	if *param == "" {
		*param = defaultValue
	}
}

func getCurrencyPair(baseParam, quoteParam string) *models.CurrencyPair {
	base := models.NewCurrency(baseParam)
	quote := models.NewCurrency(quoteParam)
	pair := models.NewCurrencyPair(base, quote)

	return &pair
}
