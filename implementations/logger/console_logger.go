package logger

import "log"

type ConsoleLogger struct{}

func (logger ConsoleLogger) Log(text string) {
	log.Println(text)
}
