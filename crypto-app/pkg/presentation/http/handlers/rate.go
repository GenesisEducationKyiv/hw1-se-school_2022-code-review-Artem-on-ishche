package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/application"
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
	logger              services.Logger
}

func (handler *RateRequestHandler) HandleRequest(ctx *gin.Context) *JSONResponse {
	var params rateRequestParameters

	if err := ctx.ShouldBind(&params); err != nil {
		handler.logger.Error("Input parameters to /rate are wrong")

		return NewJSONResponse(http.StatusBadRequest, "Input parameters are wrong")
	}

	handleEmptyParameter(&params.Base, defaultBaseName)
	handleEmptyParameter(&params.Quote, defaultQuoteName)

	pair := getCurrencyPair(params.Base, params.Quote)
	handler.logger.Debug(fmt.Sprintf("the pair after substituting empty parameters is: %s", pair.String()))

	exchangeRate, err := handler.ExchangeRateService.GetExchangeRate(*pair)
	handler.logger.Debug(fmt.Sprintf("GetExchangeRate() returned exchangeRate=%s, err=%s", exchangeRate.String(), err.Error()))

	if errors.Is(err, nil) {
		return NewJSONResponse(http.StatusOK, exchangeRate.Price)
	} else if errors.Is(err, application.ErrAPIRequestUnsuccessful) {
		return NewJSONResponse(http.StatusBadGateway, "API request has not been successful")
	} else if errors.Is(err, application.ErrAPIResponseUnmarshallError) {
		return NewJSONResponse(http.StatusBadGateway, "API returned unexpected response")
	} else {
		return NewJSONResponse(http.StatusInternalServerError, "Some unexpected error has occurred")
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
