package config

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"

	"gses2.app/api/pkg/domain/services"
)

var (
	CoinAPIKeyValue      string
	NomicsAPIKeyValue    string
	MailSlurpAPIKeyValue string

	EmailAddress  string
	EmailPassword string

	NetworkPort            string
	CryptoCurrencyProvider string
	AdminKey               string
)

func LoadEnv(loggerService services.Logger) {
	loadFile(loggerService)
	loadVariables()
}

func loadFile(loggerService services.Logger) {
	loggerService.Info("Loading .env file")

	projectDirName := "crypto-app"
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		loggerService.Error("Error loading .env file")
	}

	loggerService.Info("Successfully loaded .env file")
}

func loadVariables() {
	CoinAPIKeyValue = os.Getenv("COIN_API_KEY")
	NomicsAPIKeyValue = os.Getenv("NOMICS_API_KEY")
	MailSlurpAPIKeyValue = os.Getenv("MAILSLURP_API_KEY")

	EmailAddress = os.Getenv("EMAIL_ADDRESS")
	EmailPassword = os.Getenv("EMAIL_PASSWORD")

	NetworkPort = os.Getenv("NETWORK_PORT")
	CryptoCurrencyProvider = os.Getenv("CRYPTO_CURRENCY_PROVIDER")
	AdminKey = os.Getenv("ADMIN_KEY")
}
