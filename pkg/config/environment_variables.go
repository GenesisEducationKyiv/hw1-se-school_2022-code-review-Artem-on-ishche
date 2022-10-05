package config

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

var (
	CoinAPIKeyValue          string
	NomicsAPIKeyValue        string
	CoinMarketCapAPIKeyValue string
	MailSlurpAPIKeyValue     string

	EmailAddress  string
	EmailPassword string

	NetworkPort            string
	CryptoCurrencyProvider string
	AdminKey               string
)

func LoadEnv() {
	loadFile()
	loadVariables()
}

func loadFile() {
	projectDirName := "btc_application"
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func loadVariables() {
	CoinAPIKeyValue = os.Getenv("COIN_API_KEY")
	NomicsAPIKeyValue = os.Getenv("NOMICS_API_KEY")
	CoinMarketCapAPIKeyValue = os.Getenv("COIN_MARKETCAP_API_KEY")
	MailSlurpAPIKeyValue = os.Getenv("MAILSLURP_API_KEY")

	EmailAddress = os.Getenv("EMAIL_ADDRESS")
	EmailPassword = os.Getenv("EMAIL_PASSWORD")

	NetworkPort = os.Getenv("NETWORK_PORT")
	CryptoCurrencyProvider = os.Getenv("CRYPTO_CURRENCY_PROVIDER")
	AdminKey = os.Getenv("ADMIN_KEY")
}
