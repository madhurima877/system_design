package websocket

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*websocket.Conn),
	}
}

func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.clients[userID] = conn
}

func (h *Hub) Unregister(userId string) {
	delete(h.clients, userId)
}
func (h *Hub) SendToUser(userId string, message interface{}) error {
	conn, ok := h.clients[userId]
	if !ok {
		return nil
	}
	return conn.WriteJSON(message)
}
