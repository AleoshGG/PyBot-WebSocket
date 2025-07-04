package infrastructure

import "PyBot-WebSocket/infrastructure/controllers"

type LoadConsumers struct {
}

func NewLoadConsumers() *LoadConsumers{
	return &LoadConsumers{}
}

func (lc *LoadConsumers) Run() {
	// Declaramos los consumidores que necesitemos y los echamos a andar
	sensor_HX := controllers.NewConsumerController("sensor_HX", "hx")
	go sensor_HX.Run()

	sensor_NEO := controllers.NewConsumerController("sensor_NEO", "neo")
	go sensor_NEO.Run()
}