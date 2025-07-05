package adapters

import (
	"PyBot-WebSocket/domain/models"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	hub *models.Hub
	upgrader websocket.Upgrader
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		hub: models.NewHub(),
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {return true}},
	}
}

// Listener and Serve
func (s *WebSocketServer) ListenAndServe(addr string) error {
	http.HandleFunc("/ws/hx", s.handleWS("sensor_HX"))
	http.HandleFunc("/ws/neo", s.handleWS("sensor_NEO"))
	return http.ListenAndServe(addr, nil)
}
func (s *WebSocketServer) handleWS(sensor string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		p_id := params.Get("prototype_id")
		if p_id == "" {
			http.Error(w, "prototype_id missing", http.StatusBadRequest)
			return
		}

		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error: ", err)
			return
		}
		log.Printf("ID: %s", p_id)
		s.hub.Register(sensor, p_id, conn)
	}
}

func (s *WebSocketServer) GetHub() *models.Hub {
	return s.hub
}