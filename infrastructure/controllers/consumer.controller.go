package controllers

import (
	"PyBot-WebSocket/application/services"
	"PyBot-WebSocket/domain/models"
	"PyBot-WebSocket/infrastructure/adapters"
)

type ConsumerController struct {
	scamqp 		*services.ConsumerAMQP
	exchange 	models.Exchange
	queue    	models.Queue
	queueBind 	models.QueueBind
}

func NewConsumerController(name, key string) *ConsumerController {
	rabbit := adapters.NewRabbitMQ()
	scamqp := services.NewConsumerAMQP(rabbit)
	
	ex := models.Exchange {
		Name: "amq.topic",
		Type: "topic",
		Durable: true,
	}

	q := models.Queue {
		Name: name,
		Durable: true,
	}

	qb := models.QueueBind {
		Name: name,
		RoutingKey: key,
		Exchange: "amq.topic",
	}

	return &ConsumerController{scamqp: scamqp, exchange: ex, queue: q, queueBind: qb}
} 


func (cC *ConsumerController) Run() {
	cC.scamqp.Run(cC.exchange, cC.queue, cC.queueBind)
}