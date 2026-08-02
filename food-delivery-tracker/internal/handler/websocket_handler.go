package handler

import (
	"net/http"
	wb "system_design/food-delivery-tracker/internal/websocket"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	hub *wb.Hub
}

func NewWebSocketHandler(hub *wb.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()
	userID := "1"
	h.hub.AddConnection(userID, conn)
	defer h.hub.RemoveConnection(userID)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}

	}

}

func (h *WebSocketHandler) Notify(w http.ResponseWriter, r *http.Request) {
	userID := "1"
	message := "Order Status: Out for Delivery"

	err := h.hub.SendToUser(userID, []byte(message))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Notification sent"))
}
