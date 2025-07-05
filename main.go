package main

import (
	usecases "PyBot-WebSocket/application/useCases"
	"PyBot-WebSocket/domain/models"
	"PyBot-WebSocket/infrastructure/adapters"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// Context para control de cancelación
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1) Inicializar adaptadores
	rabbit := adapters.NewRabbitMQ()
	wsServer := adapters.NewWebSocketServer()

	// 2) Inicializar caso de uso
	chatUC := usecases.NewChatUseCase(rabbit, wsServer)

	// 3) Iniciar escucha de colas
	for _, cfg := range models.SensorConfigs() {
		go chatUC.ConsumeSensor(ctx, cfg)
	}

	// 4) Iniciar servidor WebSocket
	go func() {
		fmt.Println("WebSocket running on :8080")
		if err := wsServer.ListenAndServe(":8080"); err != nil {
			log.Fatalf("WS server error: %v", err)
		}
	}()

	// 5) Esperar señal de parada
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down...")
	cancel()
	
}
