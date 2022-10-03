package logger

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	amqpURL         = "amqp://guest:guest@localhost:5672/"
	exchangeName    = "logs"
	debugRoutingKey = "debug"
	infoRoutingKey  = "info"
	errorRoutingKey = "error"
)

type rabbitMQLogger struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQLogger() *rabbitMQLogger {
	connection, err := amqp.Dial(amqpURL)
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

func (logger *rabbitMQLogger) Debug(text string) {
	logger.sendMessage(text, debugRoutingKey)
}

func (logger *rabbitMQLogger) Info(text string) {
	logger.sendMessage(text, infoRoutingKey)
}

func (logger *rabbitMQLogger) Error(text string) {
	logger.sendMessage(text, errorRoutingKey)
}

func (logger *rabbitMQLogger) Close() {
	err := logger.channel.Close()
	handleError(err, "Failed to close a channel")

	err = logger.connection.Close()
	handleError(err, "Failed to close connection to RabbitMQ")
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
