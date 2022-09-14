package config

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

var (
	NetworkPort              string
	Filename                 string
	CoinAPIKeyValue          string
	NomicsAPIKeyValue        string
	CoinMarketCapAPIKeyValue string
	EmailAddress             string
	EmailPassword            string
	MailSlurpAPIKeyValue     string
	CryptoCurrencyProvider   string
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
	NetworkPort = os.Getenv("NETWORK_PORT")
	Filename = os.Getenv("FILENAME")
	CoinAPIKeyValue = os.Getenv("COIN_API_KEY")
	NomicsAPIKeyValue = os.Getenv("NOMICS_API_KEY")
	CoinMarketCapAPIKeyValue = os.Getenv("COIN_MARKETCAP_API_KEY")
	EmailAddress = os.Getenv("EMAIL_ADDRESS")
	EmailPassword = os.Getenv("EMAIL_PASSWORD")
	MailSlurpAPIKeyValue = os.Getenv("MAILSLURP_API_KEY")
	CryptoCurrencyProvider = os.Getenv("CRYPTO_CURRENCY_PROVIDER")
}
