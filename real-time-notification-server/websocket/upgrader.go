package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var Updgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
