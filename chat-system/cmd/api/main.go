package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/madhurima877/chat-system/api/handlers"
	"github.com/madhurima877/chat-system/internal/database"
	"github.com/madhurima877/chat-system/internal/repository"
	"github.com/madhurima877/chat-system/internal/service"
)

func main() {

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

	}
	r.Run(":8080")
}
