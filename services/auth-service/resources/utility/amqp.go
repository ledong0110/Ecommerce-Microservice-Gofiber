package utils

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	amqp_config "auth_service/config/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
	
}

func Publish(qName, text string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := amqp_config.AmqpChannel.PublishWithContext(ctx, "", qName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(text),
	})
	failOnError(err, "Failed to publish a message")
}