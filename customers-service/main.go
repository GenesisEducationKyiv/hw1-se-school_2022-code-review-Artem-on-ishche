package main

import (
	"encoding/json"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/dtm/dtmutil"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

// config
var (
	createRoute           = os.Getenv("CREATE_CUSTOMERS_ROUTE")
	dtmCoordinatorAddress = os.Getenv("DTM_COORDINATOR")
	customersServerURL    = os.Getenv("CUSTOMERS_SERVICE_URL")
	customersServerPort   = os.Getenv("CUSTOMERS_SERVICE_PORT")
	mySqlDsn              = os.Getenv("MYSQL_DSN")
)

// model
type Customer struct {
	gorm.Model
	IDTransaction string
	EmailAddress  string
	Status        string
}

// system
func main() {
	app := gin.New()

	// public order api
	app.POST(createRoute, func(c *gin.Context) {
		log.Printf("%s route called", createRoute)
		log.Println("query url was " + c.Request.URL.RawQuery)

		createCustomerRequest := struct {
			EmailAddress string `json:"emailAddress"`
		}{}

		err := c.Bind(&createCustomerRequest)
		log.Printf("binding JSON to parameters returned err={%v}", err)

		if err != nil {
			c.JSON(extractCode(err), err)
			return
		}

		globalTransactionId := dtmcli.MustGenGid(dtmCoordinatorAddress)
		req, _ := structToMap(createCustomerRequest)
		err = dtmcli.
			NewSaga(dtmCoordinatorAddress, globalTransactionId).
			Add(customersServerURL+"/register-customer", customersServerURL+"/register-customer-compensate", req).
			Submit()

		log.Printf("saga returned err = {%v}", err)

		createCustomerResponse := struct {
			Gid string `json:"gid"`
		}{Gid: globalTransactionId}

		c.JSON(extractCode(err), createCustomerResponse)
	})

	// internal order api
	app.POST("/register-customer", dtmutil.WrapHandler2(func(c *gin.Context) interface{} {
		log.Println("/register-customer route called")

		transactionId := c.Query("gid")
		log.Println("transactionId for this query is " + transactionId)

		registerCustomerRequest := struct {
			EmailAddress string `json:"email_address"`
		}{}
		err := c.BindJSON(&registerCustomerRequest)
		log.Printf("binding JSON to parameters returned err={%v}", err)
		if err != nil {
			return err
		}

		err = getDb().
			Create(&Customer{
				IDTransaction: transactionId,
				EmailAddress:  registerCustomerRequest.EmailAddress,
				Status:        "created",
			}).
			Error
		log.Printf("db create returned err={%v}", err)

		return err
	}))
	app.POST("/register-customer-compensate", dtmutil.WrapHandler2(func(c *gin.Context) interface{} {
		log.Println("/register-customer-compensate route called")

		transactionId := c.Query("gid")
		log.Println("transactionId for this query is " + transactionId)

		err := getDb().
			Model(&Customer{}).
			Where("id_transaction = ?", transactionId).
			Update("status", "canceled").
			Limit(1).
			Error
		log.Printf("db update returned err={%v}", err)

		return err
	}))

	log.Println("started")
	_ = app.Run(customersServerPort)
}

func extractCode(err error) int {
	if err == nil {
		return http.StatusOK
	} else {
		return http.StatusInternalServerError
	}
}

func getDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(mySqlDsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&Customer{})

	return db
}

func structToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap)
	return
}
