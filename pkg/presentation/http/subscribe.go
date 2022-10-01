package http

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

func (handler SubscribeRequestHandler) GetPath() string {
	return "/subscribe"
}

func (handler SubscribeRequestHandler) GetMethod() string {
	return "POST"
}

func (handler SubscribeRequestHandler) HandleRequest(ctx *gin.Context) {
	var params subscribeRequestParameters

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, "Input parameters are wrong")

		return
	}

	emailAddress, err := models.NewEmailAddress(params.EmailAddrString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Provided email address is wrong")

		return
	}

	pair := getCurrencyPair(params.Base, params.Quote)

	err = handler.RateSubscriptionService.Subscribe(emailAddress, pair)
	if err == nil {
		ctx.JSON(http.StatusOK, "Success")
	} else if isEmailAlreadySaved(err, emailAddress.String()) {
		ctx.JSON(http.StatusConflict, "This email address is already saved")
	} else {
		ctx.JSON(http.StatusInternalServerError, "Error when saving the email address")
	}
}

func isEmailAlreadySaved(err error, emailAddressString string) bool {
	return err.Error() == application.ErrEmailAddressAlreadyExists(emailAddressString).Error()
}
