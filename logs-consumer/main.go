package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	amqpURL      = "amqp://guest:guest@localhost:5672/"
	exchangeName = "logs"
	queueName    = "logs-queue"
	bindingKey   = "error"
)

func main() {
	connection := getConnection()
	defer connection.Close()

	channel := getChannelWithQueueSetUp(connection)
	defer channel.Close()

	messages := getMessages(channel)
	readForever(messages)
}

func getConnection() *amqp.Connection {
	connection, err := amqp.Dial(amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")

	return connection
}

func getChannelWithQueueSetUp(conn *amqp.Connection) *amqp.Channel {
	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = channel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	_, err = channel.QueueDeclare(
		queueName,
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = channel.QueueBind(
		queueName,
		bindingKey,
		exchangeName,
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	return channel
}

func getMessages(channel *amqp.Channel) <-chan amqp.Delivery {
	messages, err := channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	return messages
}

func readForever(messages <-chan amqp.Delivery) {
	var forever chan struct{}

	go func() {
		for d := range messages {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
