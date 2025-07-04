package adapters

import (
	"PyBot-WebSocket/domain/models"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
	conn, err := amqp.Dial(os.Getenv("URL_RABBIT"))
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")


	return &RabbitMQ{conn: conn, ch: ch}  
}

func (r *RabbitMQ) ConsumeQueue(ex models.Exchange, q models.Queue, qb models.QueueBind) {
	err := r.ch.ExchangeDeclare(
        ex.Name,   // name
        ex.Type, // type
        ex.Durable,     // durable
        false,    // auto-deleted
        false,    // internal
        false,    // no-wait
        nil,      // arguments
    )
    failOnError(err, "Failed to declare an exchange")

    que, err := r.ch.QueueDeclare(
        q.Name,    // name
        q.Durable, // durable
        false, // delete when unused
        false,  // exclusive
        false, // no-wait
        nil,   // arguments
    )
    failOnError(err, "Failed to declare a queue")

    err = r.ch.QueueBind(
        qb.Name, // queue name
        qb.RoutingKey,     // routing key
        qb.Exchange, // exchange
        false,
        nil,
    )
    failOnError(err, "Failed to bind a queue")

    msgs, err := r.ch.Consume(
        que.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
    	nil,    // args
    )
    failOnError(err, "Failed to register a consumer")

    var forever chan struct{}

    go func() {
        for d := range msgs {
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