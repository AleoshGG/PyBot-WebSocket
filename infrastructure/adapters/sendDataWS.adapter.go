package adapters

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type GorillaClient struct {
	conn *websocket.Conn
}

func NewGorillaClient(room string, client string) *GorillaClient {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	q := u.Query()
	q.Set("room", room)
	q.Set("client", client)
	u.RawQuery = q.Encode()

	log.Printf("Conectando a %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error al conectar: ", err)
	}

	return &GorillaClient{conn: conn}
}

func (gc *GorillaClient) SendData(data []byte) {
	err := gc.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("Hubo un error: ", err)
		return
	}
}

func (gc *GorillaClient) Close() error {
    return gc.conn.Close()
}
