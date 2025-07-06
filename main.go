package main

import (
	usecases "PyBot-WebSocket/application/useCases"
	"PyBot-WebSocket/domain/models"
	"PyBot-WebSocket/infrastructure/adapters"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	godotenv.Load()

	// Context para control de cancelación
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1) Inicializar adaptadores
	rabbit := adapters.NewRabbitMQ()
	wsServer := adapters.NewGorilla()

	// 2) Inicializar caso de uso
	chatUC := usecases.NewChatUseCase(rabbit, wsServer)

	// 3) Iniciar escucha de colas
	for _, cfg := range models.SensorConfigs() {
		go chatUC.ConsumeSensor(ctx, cfg)
	}

	// 4) Iniciar servidor WebSocket
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/hx",  wsServer.HandleWS("sensor_HX"))
	mux.HandleFunc("/ws/neo", wsServer.HandleWS("sensor_NEO"))
	mux.HandleFunc("/ws/cam", wsServer.HandleWS("sensor_CAM"))

	// Configuración de CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   adapters.AllowedOrigins,
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
        // MaxAge: 600, // opcional: cachear preflight en segundos
    })

	server := &http.Server{
		Addr: ":8080",
		Handler: c.Handler(mux),
	}

	go func() {
		fmt.Println("WebSocket running on :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
