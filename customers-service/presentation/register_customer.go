package presentation

import (
	"log"

	"github.com/gin-gonic/gin"

	"customers-service/model"
)

func registerCustomerHandler(c *gin.Context) interface{} {
	log.Println("/register-customer route called")

	transactionId := c.Query("gid")
	log.Println("transactionId for this query is " + transactionId)

	var registerCustomerRequest createCustomerRequestType
	err := c.BindJSON(&registerCustomerRequest)

	log.Printf("binding JSON to parameters returned err={%v}", err)
	if err != nil {
		return err
	}

	err = getDb().
		Create(&model.Customer{
			IDTransaction: transactionId,
			EmailAddress:  registerCustomerRequest.EmailAddress,
			Status:        "created",
		}).
		Error
	log.Printf("db create returned err={%v}", err)

	return err
}
