package presentation

import (
	"customers-service/config"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	app := gin.New()

	// public order api
	app.POST(config.CreateRoute, createHandler)

	// internal order api
	app.POST("/register-customer", dtmutil.WrapHandler2(registerCustomerHandler))
	app.POST("/register-customer-compensate", dtmutil.WrapHandler2(registerCustomerCompensateHandler))

	return app
}
