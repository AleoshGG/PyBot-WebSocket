package infrastructure

import (
	"PyBot-WebSocket/application/services"
	"PyBot-WebSocket/domain/models"
	"PyBot-WebSocket/infrastructure/adapters"
)

type SendData struct {
	scamqp 		*services.ConsumerAMQP
	exchange 	models.Exchange
	queue    	models.Queue
	queueBind 	models.QueueBind
}

func NewSendData() *SendData {
	rabbit := adapters.NewRabbitMQ()
	scamqp := services.NewConsumerAMQP(rabbit)
	
	ex := models.Exchange {
		Name: "amq.topic",
		Type: "topic",
		Durable: true,
	}

	q := models.Queue {
		Name: "sensor_HX",
		Durable: true,
	}

	qb := models.QueueBind {
		Name: "sensor_HX",
		RoutingKey: "hx",
		Exchange: "amq.topic",
	}

	return &SendData{scamqp: scamqp, exchange: ex, queue: q, queueBind: qb}
} 


func (sd *SendData) Run() {
	sd.scamqp.Run(sd.exchange, sd.queue, sd.queueBind)
}