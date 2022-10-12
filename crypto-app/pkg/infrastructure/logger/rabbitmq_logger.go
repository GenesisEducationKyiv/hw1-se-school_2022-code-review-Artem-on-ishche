package logger

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"gses2.app/api/pkg/config"
)

const (
	exchangeName    = "logs"
	debugRoutingKey = "debug"
	infoRoutingKey  = "info"
	errorRoutingKey = "error"
	timeFormat      = "2006-01-02 15:04:05 "
)

type rabbitMQLogger struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQLogger() *rabbitMQLogger {
	connection, err := amqp.Dial(config.AmqpUrl)
	handleError(err, "Failed to connect to RabbitMQ")

	channel, err := connection.Channel()
	handleError(err, "Failed to open a channel")

	err = channel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Failed to declare an exchange")

	return &rabbitMQLogger{
		connection: connection,
		channel:    channel,
	}
}

func (logger *rabbitMQLogger) Debug(messageText string) {
	logger.sendMessage(getTimeText()+messageText, debugRoutingKey)
}

func (logger *rabbitMQLogger) Info(messageText string) {
	logger.sendMessage(getTimeText()+messageText, infoRoutingKey)
}

func (logger *rabbitMQLogger) Error(messageText string) {
	logger.sendMessage(getTimeText()+messageText, errorRoutingKey)
}

func (logger *rabbitMQLogger) Close() {
	err := logger.channel.Close()
	handleError(err, "Failed to close a channel")

	err = logger.connection.Close()
	handleError(err, "Failed to close connection to RabbitMQ")
}

func getTimeText() string {
	return time.Now().Format(timeFormat)
}

func (logger *rabbitMQLogger) sendMessage(text, key string) {
	err := logger.channel.PublishWithContext(context.Background(),
		exchangeName,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(text),
		},
	)
	handleError(err, "Failed to publish a message")
}

func handleError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
