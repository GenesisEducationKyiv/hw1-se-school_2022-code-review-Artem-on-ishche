package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	CoinAPIKeyValue      string
	NomicsAPIKeyValue    string
	MailSlurpAPIKeyValue string

	EmailAddress  string
	EmailPassword string

	CustomersServiceUrl string
	CreateCustomerRoute string

	NetworkPort            string
	CryptoCurrencyProvider string
	AdminKey               string

	AmqpUrl string
)

func LoadEnv() {
	loadFile()
	loadVariables()
}

func loadFile() {
	log.Println("Loading .env file")

	maxDirectoryDepth := 6
	for i := 0; i < maxDirectoryDepth; i++ {
		escapeSequence := strings.Repeat("../", i)

		err := godotenv.Load("./" + escapeSequence + ".env")
		if err == nil {
			log.Println("Success")

			return
		}
	}

	log.Fatal("Failed to load .env file")
}

func loadVariables() {
	CoinAPIKeyValue = os.Getenv("COIN_API_KEY")
	NomicsAPIKeyValue = os.Getenv("NOMICS_API_KEY")
	MailSlurpAPIKeyValue = os.Getenv("MAILSLURP_API_KEY")

	EmailAddress = os.Getenv("EMAIL_ADDRESS")
	EmailPassword = os.Getenv("EMAIL_PASSWORD")

	CustomersServiceUrl = os.Getenv("CUSTOMERS_SERVICE_URL")
	CreateCustomerRoute = os.Getenv("CREATE_CUSTOMERS_ROUTE")

	NetworkPort = os.Getenv("NETWORK_PORT")
	CryptoCurrencyProvider = os.Getenv("CRYPTO_CURRENCY_PROVIDER")
	AdminKey = os.Getenv("ADMIN_KEY")

	AmqpUrl = os.Getenv("AMQP_URL")
}
