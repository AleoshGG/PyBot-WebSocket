package usecases

import (
	"PyBot-WebSocket/domain/models"
	"PyBot-WebSocket/infrastructure/adapters"
	"context"
	"encoding/json"
	"log"
)

// ChatUseCase coordina RabbitMQ y WebSocket
type ChatUseCase struct {
	rmq *adapters.RabbitMQ
	hub *models.Hub
}

func NewChatUseCase(rmq *adapters.RabbitMQ, wsServer *adapters.WebSocketServer) *ChatUseCase {
	return &ChatUseCase{rmq: rmq, hub: wsServer.GetHub()}
}

func (uc *ChatUseCase) ConsumeSensor(ctx context.Context, cfg models.SensorConfig) {
	msgs, err := uc.rmq.Consume(cfg.Exchange, cfg.Queue, cfg.RoutingKey)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case d := <-msgs:
			// Extraer prototype_id
			var pid string
			switch cfg.Queue {
			case "sensor_HX":
				var data models.HX711
				json.Unmarshal(d.Body, &data)
				pid = data.Prototype_id
				log.Printf("Data [hx]: %s", d.Body)
			case "sensor_NEO":
				var data models.GPS
				json.Unmarshal(d.Body, &data)
				pid = data.Prototype_id
				log.Printf("Data [neo]: %s", d.Body)
			case "sensor_ABC":
				// ...
			default: log.Printf("Data: %s", d.Body)
			}
			log.Printf("ID: %s", pid)
			log.Printf("queue: %s", cfg.Queue)
			log.Printf("body: %s", d.Body)
			uc.hub.Send(models.SensorMessage{Sensor: cfg.Queue, Prototype_id: pid, Payload: d.Body})
		}
		}
}
