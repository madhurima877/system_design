package handlers

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/madhurima877/chat-system/api/websocket"
	"github.com/madhurima877/chat-system/internal/models"
	"github.com/madhurima877/chat-system/proto/chat"
)

type ChatHandler struct {
	client chat.ChatServiceClient
	hub    *websocket.Hub
}

func NewChatHandler(client chat.ChatServiceClient, hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{
		client: client,
		hub:    hub,
	}
}
func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req models.MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	senderID := c.GetString("userID")
	resp, err := h.client.SendMessage(context.Background(), &chat.SendMessageRequest{
		SenderId:   senderID,
		ReceiverId: req.ReceiverID,
		Content:    req.Content,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := h.hub.SendToUser(req.ReceiverID, gin.H{
		"sender_id": senderID,
		"content":   req.Content,
	}); err != nil {
		log.Println("websocket:", err)
	}
	c.JSON(200, resp)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	senderID := c.GetString("userID")
	receiverID := c.Param("receiverId")
	resp, err := h.client.GetMessages(context.Background(), &chat.GetMessagesRequest{
		SenderId:   senderID,
		ReceiverId: receiverID,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}
