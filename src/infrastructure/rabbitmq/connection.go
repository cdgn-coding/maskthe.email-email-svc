package rabbitmq

import (
	"email-svc/src/infrastructure/configuration"
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateConnection(logger configuration.Logger, url string) *amqp.Channel {
	connection, err := amqp.Dial(url)

	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ %v", err)
		panic(err)
	}

	channel, err := connection.Channel()
	if err != nil {
		logger.Fatal("Failed to open a channel %v", err)
		panic(err)
	}

	channelErrors := channel.NotifyClose(make(chan *amqp.Error, 10))
	go func() {
		for err := range channelErrors {
			if err != nil {
				logger.Fatal("Channel error: %v", err)
			}
		}
	}()

	return channel
}
