package repositories

import amqp "github.com/rabbitmq/amqp091-go"

type ConsumerAMQP interface {
	Consume(exchange, queue, key string) (<-chan amqp.Delivery, error)
}