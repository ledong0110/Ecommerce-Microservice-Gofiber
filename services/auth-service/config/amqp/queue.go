package amqp_conn

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

var AmqpChannel *amqp.Channel
var EmailService amqp.Queue