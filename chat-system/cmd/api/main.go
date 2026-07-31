package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/madhurima877/chat-system/api/handlers"
	"github.com/madhurima877/chat-system/api/middleware"
	"github.com/madhurima877/chat-system/api/websocket"
	"github.com/madhurima877/chat-system/internal/database"
	"github.com/madhurima877/chat-system/internal/repository"
	"github.com/madhurima877/chat-system/internal/service"
	"github.com/madhurima877/chat-system/proto/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50059", grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatalln(err)

	}
	defer conn.Close()
	chatClient := chat.NewChatServiceClient(conn)
	hub := websocket.NewHub()
	wsnHandler := websocket.NewHandler(hub)
	chatHandler := handlers.NewChatHandler(chatClient, hub)

	client, err := database.ConnectMongo()
	if err != nil {
		log.Fatalln(err)

	}

	userRepo := repository.NewUserRepository(client)
	userService := service.NewAuthService(userRepo)
	userHandler := handlers.NewAuthHandler(userService)
	fmt.Println("MongoDB Connected")
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.POST("/send/message", middleware.AuthMiddleware(), chatHandler.SendMessage)
		v1.GET("/messages/:receiverId", middleware.AuthMiddleware(), chatHandler.GetMessages)
		v1.GET("/ws", middleware.AuthMiddleware(), wsnHandler.Connect)

	}
	if err := r.Run(":8080"); err != nil {
		log.Fatalln(err)
	}

}
