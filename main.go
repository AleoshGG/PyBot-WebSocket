package main

import (
	usecases "PyBot-WebSocket/application/useCases"
	"PyBot-WebSocket/infrastructure"
	"PyBot-WebSocket/infrastructure/controllers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// Canal para capturar errores del servidor HTTP
	serverErrors := make(chan error, 1)

	// Levantaar el servidor ws en un goroutine
	go func() {
		chatService := usecases.NewChat()
		wsHandler := controllers.NewWebSocketController(chatService)

		mux := http.NewServeMux()
        mux.HandleFunc("/ws", wsHandler.HandleWS)

        server := &http.Server{
            Addr:    ":8080",
            Handler: mux,
        }

        fmt.Println("Servidor WebSocket corriendo en :8080")

		// Envia al canal si ListenAndServe falla
        serverErrors <- server.ListenAndServe()
	} ()

	// Correr consumidores
	consumers := infrastructure.NewLoadConsumers()
	consumers.Run()
	
	// Captura de señal para graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)

    // Bloquea esperando o bien un error de servidor, o Ctrl+C
    select {
    case err := <-serverErrors:
        log.Fatalf("El servidor WS terminó con error: %v", err)
    case <-stop:
        log.Println("Recibida señal de parada, apagando...")

        // Aquí podrías hacer shutdown controlado de tu WS server si lo almacenaste
        // en una variable accesible, o notificar a los consumidores que detengan.
        // Ejemplo de shutdown en el mismo código:
        _, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        // server.Shutdown(ctx)  // si guardaste `server` en un scope externo
    }

    log.Println("Aplicación finalizada.")
	
}
