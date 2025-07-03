package controllers

import (
	usecases "PyBot-WebSocket/application/useCases"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {return true},
}

type WebSocketController struct {
	chat *usecases.Chat
}

func NewWebSocketController(chat *usecases.Chat) *WebSocketController {
	return &WebSocketController{chat: chat}
}

func (wsc *WebSocketController) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error: ", err)
		return
	}
	defer conn.Close()

	roomID := r.URL.Query().Get("room")
	clientID := r.URL.Query().Get("client")

	client := wsc.chat.JoinRoom(roomID, clientID)

	go func() {
		for msg := range client.Send {
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	} ()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error: ", err)
			break
		}
		fmt.Printf("Received from %s: %s\n", clientID, msg)
		wsc.chat.SendMessage(roomID, msg)
	}

}