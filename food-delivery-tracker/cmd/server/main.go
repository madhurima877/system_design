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
	service := service.NewUserService(repo, hub)
	wsnHandler := handler.NewWebSocketHandler(hub)
	handler := handler.NewUserHandler(service)
	grpcHandler := grpc.NewUserServer(service)

	user.RegisterUserServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)

	http.HandleFunc("/create/user", handler.CreateUser)
	http.HandleFunc("/ws", wsnHandler.HandleWebSocket)
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
	grpcServer.GracefulStop()

}
