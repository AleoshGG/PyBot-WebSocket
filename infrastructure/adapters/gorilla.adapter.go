package adapters

import (
	"PyBot-WebSocket/domain/models"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Gorilla struct {
	hub *models.Hub
	upgrader websocket.Upgrader
}

func NewGorilla() *Gorilla {
	return &Gorilla{
		hub: models.NewHub(),
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {return true}},
	}
}

// Listener and Serve
func (s *Gorilla) ListenAndServe(addr string) error {
	http.HandleFunc("/ws/hx", s.HandleWS("sensor_HX"))
	http.HandleFunc("/ws/neo", s.HandleWS("sensor_NEO"))
	http.HandleFunc("/ws/cam", s.HandleWS("sensor_CAM"))
	return http.ListenAndServe(addr, nil)
}

func (s *Gorilla) HandleWS(sensor string) http.HandlerFunc {
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

func (s *Gorilla) GetHub() *models.Hub {
	return s.hub
}