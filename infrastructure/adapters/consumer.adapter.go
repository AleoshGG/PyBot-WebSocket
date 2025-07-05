package adapters

import (
	"PyBot-WebSocket/domain/models"
	"context"
	"encoding/json"
	"errors"
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

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    if err := processData(msgs, ctx, q.Name); err != nil {
        log.Printf("Error: %s", err)
    }
}

func processData(msgs <-chan amqp.Delivery, ctx context.Context, qname string) error {
    wsClients := make(map[string]*GorillaClient)

    for {
        select {
            case <- ctx.Done():
                for _, client := range wsClients {
                    client.Close()
                }
                return ctx.Err()

            case d, ok := <- msgs:
                if !ok {
                    return errors.New("Message channel closes")
                }

                var prototype_id string
                switch qname {
                case "sensor_HX":
                    var data models.HX711
                    if err := json.Unmarshal(d.Body, &data); err != nil {
                        log.Printf("Error parsing sensor_HX data: %s", err)
                        continue
                    }
                    prototype_id = data.Prototype_id
                case "sensor_NEO":
                    var data models.GPS
                    if err := json.Unmarshal(d.Body, &data); err != nil {
                        log.Printf("Error parsing sensor_NEO data: %s", err)
                        continue
                    }
                    prototype_id = data.Prototype_id
                default: 
                    log.Printf("No handler for queue %s: %s", qname, d.Body)
                    continue
                }
                
                client, exists := wsClients[prototype_id]
                if !exists {
                    client = NewGorillaClient(qname, prototype_id)
                    wsClients[prototype_id] = client
                }

                client.SendData(d.Body)
        }   
    }
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
}