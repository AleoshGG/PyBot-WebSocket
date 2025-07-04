package repositories

import "PyBot-WebSocket/domain/models"

type IConsumerAMQP interface {
	ConsumeQueue(ex models.Exchange, q models.Queue, qb models.QueueBind)
}