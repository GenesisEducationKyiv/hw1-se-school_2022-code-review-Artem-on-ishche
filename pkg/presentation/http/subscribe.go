package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

type subscribeRequestParameters struct {
	EmailAddrString string `form:"email" binding:"required"`
}

type SubscribeRequestHandler struct {
	AddEmailAddressService application.AddEmailAddressService
}

func (handler SubscribeRequestHandler) GetPath() string {
	return "/subscribe"
}

func (handler SubscribeRequestHandler) GetMethod() string {
	return "POST"
}

func (handler SubscribeRequestHandler) HandleRequest(c *gin.Context) {
	var params subscribeRequestParameters

	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Required parameter 'email' is missing")

		return
	}

	emailAddress, err := models.NewEmailAddress(params.EmailAddrString)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Provided email address is wrong")

		return
	}

	err = handler.AddEmailAddressService.AddEmailAddress(*emailAddress)
	if err == nil {
		c.JSON(http.StatusOK, "Success")
	} else if isEmailAlreadySaved(err, emailAddress.String()) {
		c.JSON(http.StatusConflict, "This email address is already saved")
	} else {
		c.JSON(http.StatusInternalServerError, "Error when saving the email address")
	}
}

func isEmailAlreadySaved(err error, emailAddressString string) bool {
	return err.Error() == services.ErrEmailAddressAlreadyExists(emailAddressString).Error()
}
