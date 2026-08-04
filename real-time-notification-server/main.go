package main

import (
	"fmt"
	"net/http"
	"system_design/real-time-notification-server/websocket"
	"system_design/real-time-notification-server/workers"

	"github.com/gin-gonic/gin"
)

func main() {
	hub := websocket.NewHub()
	pool := workers.NewWorkerPool(hub)
	pool.Start()

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Real-Time Notification Server"})

		})
		v1.POST("/notify", func(c *gin.Context) {
			var event workers.Job
			if err := c.ShouldBindJSON(&event); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			pool.Jobs <- event

			c.JSON(200, gin.H{"message": "Notification Received"})
		})
		v1.GET("/ws", func(c *gin.Context) {
			conn, err := websocket.Updgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			client := &websocket.Client{
				Conn: conn,
			}
			hub.Clients[client] = true
			fmt.Println("Client connected:", len(hub.Clients))
			defer conn.Close()
			defer func() {
				delete(hub.Clients, client)
			}()
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}

		})
	}
	router.Run(":8080")
}
