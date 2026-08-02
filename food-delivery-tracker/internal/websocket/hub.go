package websocket

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	connections map[string]*websocket.Conn
	mu          sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[string]*websocket.Conn),
	}
}
func (h *Hub) AddConnection(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.connections[userID] = conn
}

func (h *Hub) RemoveConnection(userId string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.connections, userId)
}
func (h *Hub) GetConnection(userID string) (*websocket.Conn, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	val, ok := h.connections[userID]
	if !ok {
		return nil, false
	}
	return val, true
}
func (h *Hub) SendToUser(userID string, message []byte) error {
	conn, ok := h.GetConnection(userID)
	if !ok {
		return errors.New("no connection found")
	}
	return conn.WriteMessage(websocket.TextMessage, message)
}
