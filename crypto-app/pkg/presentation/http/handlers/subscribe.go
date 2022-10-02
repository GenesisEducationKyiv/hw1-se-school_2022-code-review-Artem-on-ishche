package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type subscribeRequestParameters struct {
	EmailAddrString string `form:"email" binding:"required"`
	Base            string `form:"base"`
	Quote           string `form:"quote"`
}

type SubscribeRequestHandler struct {
	RateSubscriptionService application.RateSubscriptionService
	logger                  services.Logger
}

func (handler SubscribeRequestHandler) HandleRequest(ctx *gin.Context) *JSONResponse {
	var params subscribeRequestParameters

	if err := ctx.ShouldBind(&params); err != nil {
		handler.logger.Error("Input parameters to /subscribe are wrong")

		return NewJSONResponse(http.StatusBadRequest, "Input parameters are wrong")
	}

	emailAddress, err := models.NewEmailAddress(params.EmailAddrString)
	if err != nil {
		handler.logger.Error("Email parameter to /subscribe is invalid: " + params.EmailAddrString)

		return NewJSONResponse(http.StatusBadRequest, "Provided email address is invalid")
	}

	pair := getCurrencyPair(params.Base, params.Quote)

	err = handler.RateSubscriptionService.Subscribe(emailAddress, pair)
	handler.logger.Debug(fmt.Sprintf("SendRateEmails() returned err=%s", err.Error()))

	if err == nil {
		return NewJSONResponse(http.StatusOK, "Success")
	} else if isEmailAlreadySaved(err, emailAddress.String()) {
		return NewJSONResponse(http.StatusConflict, "This email address is already saved")
	} else {
		return NewJSONResponse(http.StatusInternalServerError, "Error when saving the email address")
	}
}

func isEmailAlreadySaved(err error, emailAddressString string) bool {
	return err.Error() == application.ErrEmailAddressAlreadyExists(emailAddressString).Error()
}
