package logger

import "log"

type ConsoleLogger struct{}

func (logger ConsoleLogger) Debug(text string) {
	log.Println("[DEBUG] " + text)
}

func (logger ConsoleLogger) Info(text string) {
	log.Println("[INFO] " + text)
}

func (logger ConsoleLogger) Error(text string) {
	log.Println("[ERROR] " + text)
}
