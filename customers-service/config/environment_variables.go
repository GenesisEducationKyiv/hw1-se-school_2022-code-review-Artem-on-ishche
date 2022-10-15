package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	CustomersServerURL  string
	CustomersServerPort string
	CreateRoute         string

	DtmCoordinatorAddress string
	MySqlDsn              string
)

func LoadEnv() {
	loadFile()
	loadVariables()
}

func loadFile() {
	log.Println("Loading .env file")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	log.Println("Success")
}

func loadVariables() {
	CustomersServerURL = os.Getenv("CUSTOMERS_SERVICE_URL")
	CustomersServerPort = os.Getenv("CUSTOMERS_SERVICE_PORT")
	CreateRoute = os.Getenv("CREATE_CUSTOMERS_ROUTE")

	DtmCoordinatorAddress = os.Getenv("DTM_COORDINATOR")
	MySqlDsn = os.Getenv("MYSQL_DSN")
}
