package presentation

import (
	"log"

	"github.com/gin-gonic/gin"

	"customers-service/model"
)

func registerCustomerCompensateHandler(c *gin.Context) interface{} {
	log.Println("/register-customer-compensate route called")

	transactionId := c.Query("gid")
	log.Println("transactionId for this query is " + transactionId)

	err := getDb().
		Model(&model.Customer{}).
		Where("id_transaction = ?", transactionId).
		Update("status", "canceled").
		Limit(1).
		Error
	log.Printf("db update returned err={%v}", err)

	return err
}
