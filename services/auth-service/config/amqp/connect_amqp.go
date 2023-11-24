package amqp_conn

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}

}

func ConnectAmqp(uri string) *amqp.Connection {
	conn, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	AmqpChannel, err = conn.Channel()

	if err != nil {
        panic(err)
    }

	
	EmailService, err = AmqpChannel.QueueDeclare(
		"EmailService",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	log.Println("RabitMQ connected !")
	return conn

}
