package presentation

import (
	"log"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/gin-gonic/gin"

	"customers-service/config"
)

type createCustomerRequestType struct {
	EmailAddress string `json:"emailAddress"`
}

type createCustomerResponseType struct {
	Gid string `json:"gid"`
}

func createHandler(c *gin.Context) {
	log.Printf("%s route called", config.CreateRoute)
	log.Println("query url was " + c.Request.URL.RawQuery)

	var createCustomerRequest createCustomerRequestType
	err := c.Bind(&createCustomerRequest)
	log.Printf("binding JSON to parameters returned err={%v}", err)

	if err != nil {
		c.JSON(extractCode(err), err)
		return
	}

	globalTransactionId := dtmcli.MustGenGid(config.DtmCoordinatorAddress)
	req, _ := structToMap(createCustomerRequest)
	err = dtmcli.
		NewSaga(config.DtmCoordinatorAddress, globalTransactionId).
		Add(config.CustomersServerURL+"/register-customer", config.CustomersServerURL+"/register-customer-compensate", req).
		Submit()

	log.Printf("saga returned err = {%v}", err)

	createCustomerResponse := createCustomerResponseType{Gid: globalTransactionId}
	c.JSON(extractCode(err), createCustomerResponse)
}
