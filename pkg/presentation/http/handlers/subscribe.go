package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
)

type subscribeRequestParameters struct {
	EmailAddrString string `form:"email" binding:"required"`
	Base            string `form:"base"`
	Quote           string `form:"quote"`
}

type SubscribeRequestHandler struct {
	RateSubscriptionService application.RateSubscriptionService
}

func (handler SubscribeRequestHandler) HandleRequest(ctx *gin.Context) *JSONResponse {
	var params subscribeRequestParameters

	if err := ctx.ShouldBind(&params); err != nil {
		return NewJSONResponse(http.StatusBadRequest, "Input parameters are wrong")
	}

	emailAddress, err := models.NewEmailAddress(params.EmailAddrString)
	if err != nil {
		return NewJSONResponse(http.StatusBadRequest, "Provided email address is wrong")
	}

	pair := getCurrencyPair(params.Base, params.Quote)

	err = handler.RateSubscriptionService.Subscribe(emailAddress, pair)
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
