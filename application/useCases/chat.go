package usecases

import "PyBot-WebSocket/domain/models"

type Chat struct {
	Rooms map[string]*models.Room
}

func NewChat() *Chat {
	return &Chat{Rooms: map[string]*models.Room{}}
}

func (c *Chat) JoinRoom(roomID, clientID string) *models.Client {
	room, exists := c.Rooms[roomID]
	if !exists {
		room = models.NewRoom(roomID)
		c.Rooms[roomID] = room
	}

	client := &models.Client{
		ID: clientID,
		Send: make(chan []byte),
	}

	room.AddClient(client)
	return client
}

func (c *Chat) SendMessage(roomID string, msg []byte) {
	if room, ok := c.Rooms[roomID]; ok {
		room.Broadcast(msg)
	}
}
