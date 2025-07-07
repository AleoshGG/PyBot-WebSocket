package repositories

import (
	"PyBot-WebSocket/domain/models"
	"net/http"
)

type WebSocketServer interface {
	// ListenAndServe(addr string) error
	HandleWS(sensor string) http.HandlerFunc
	GetHub() *models.Hub
}