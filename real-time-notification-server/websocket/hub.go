package websocket

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	Clients map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[*Client]bool),
	}
}
func (h *Hub) Broadcast(message []byte) {
	for client := range h.Clients {
		client.Mu.Lock()
		client.Conn.WriteMessage(websocket.TextMessage, message)
		client.Mu.Unlock()
	}
}
