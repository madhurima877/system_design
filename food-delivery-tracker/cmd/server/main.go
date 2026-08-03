package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"system_design/food-delivery-tracker/internal/db"
	"system_design/food-delivery-tracker/internal/grpc"
	"system_design/food-delivery-tracker/internal/notification"
	"system_design/food-delivery-tracker/internal/websocket"
	"system_design/food-delivery-tracker/proto/user"

	"system_design/food-delivery-tracker/internal/handler"
	"system_design/food-delivery-tracker/internal/repository"
	"system_design/food-delivery-tracker/internal/service"

	grpcn "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	lis, err := net.Listen("tcp", ":50060")
	if err != nil {
		panic(err)
	}

	grpcServer := grpcn.NewServer()

	log.Println("Database Connected")
	repo := repository.NewUserRepository(conn)
	hub := websocket.NewHub()
	userService := service.NewUserService(repo, hub)
	wsnHandler := handler.NewWebSocketHandler(hub)
	userHandler := handler.NewUserHandler(userService)
	grpcHandler := grpc.NewUserServer(userService)

	user.RegisterUserServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)
	notificationService := notification.NewService(hub)
	orderRepo := repository.NewOrderRepository(conn)
	orderService := service.NewOrderService(orderRepo, hub, notificationService)
	orderHandler := handler.NewOrderHandler(orderService)

	notificationService.Start()

	http.HandleFunc("/create/user", userHandler.CreateUser)
	http.HandleFunc("/ws", wsnHandler.HandleWebSocket)
	http.HandleFunc("/order/status", orderHandler.UpdateOrderStatus)

	http.HandleFunc("/notify", wsnHandler.Notify)
	go func() {
		log.Println("HTTP Server started on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		log.Println("gRPC Server started on :50060")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	notificationService.Stop()
	grpcServer.GracefulStop()

}
