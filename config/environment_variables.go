package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	NetworkPort   string
	Filename      string
	APIKeyValue   string
	EmailAddress  string
	EmailPassword string
)

func LoadEnv() {
	loadFile()
	loadVariables()
}

func loadFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
}

func loadVariables() {
	NetworkPort = os.Getenv("NETWORK_PORT")
	Filename = os.Getenv("FILENAME")
	APIKeyValue = os.Getenv("API_KEY")
	EmailAddress = os.Getenv("EMAIL_ADDRESS")
	EmailPassword = os.Getenv("EMAIL_PASSWORD")
}
