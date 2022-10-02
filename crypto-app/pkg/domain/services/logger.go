package services

type Logger interface {
	Debug(text string)
	Info(text string)
	Error(text string)
}
