package models

type Room struct {
	ID      string
	Clients map[string]*Client
}

func NewRoom(id string) *Room {
	return &Room{ID: id, Clients: make(map[string]*Client)}
}

func (r *Room) AddClient(c *Client) {
	r.Clients[c.ID] = c
}

func (r *Room) Broadcast(message []byte) {
	for _, c := range r.Clients {
		c.Send <- message
	}
}

