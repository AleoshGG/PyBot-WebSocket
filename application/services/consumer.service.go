package services

import (
	"PyBot-WebSocket/domain/models"
	"PyBot-WebSocket/domain/repositories"
)

type ConsumerAMQP struct {
	cq repositories.IConsumerAMQP
}

func NewConsumerAMQP(cq repositories.IConsumerAMQP) *ConsumerAMQP {
	return &ConsumerAMQP{cq: cq}
}

func (c *ConsumerAMQP) Run(ex models.Exchange, q models.Queue, qb models.QueueBind) {
	c.cq.ConsumeQueue(ex, q, qb)
}