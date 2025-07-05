package adapters

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ encapsula conexión y canal
type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
	conn, err := amqp.Dial(os.Getenv("URL_RABBIT"))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return &RabbitMQ{conn: conn, ch: ch}
}

func (r *RabbitMQ) Consume(exchange, queue, key string) (<-chan amqp.Delivery, error) {
	err := r.ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	_, err = r.ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	err = r.ch.QueueBind(queue, key, exchange, false, nil)
	if err != nil {
		return nil, err
	}
	return r.ch.Consume(queue, "", true, false, false, false, nil)
}