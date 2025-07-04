package main

import (
	usecases "PyBot-WebSocket/application/useCases"
	"PyBot-WebSocket/infrastructure"
	"PyBot-WebSocket/infrastructure/controllers"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// Iniciar con la escucha
	consumer := infrastructure.NewSendData()
	consumer.Run()

	chatService := usecases.NewChat()
	wsHandler := controllers.NewWebSocketController(chatService)

	http.HandleFunc("/ws", wsHandler.HandleWS)

	fmt.Println("Servidor WebSocket corriendo en :8080")
	http.ListenAndServe(":8080", nil)

	
}
