package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/madhurima877/chat-system/internal/database"
	chatgrpc "github.com/madhurima877/chat-system/internal/grpc"
	"github.com/madhurima877/chat-system/internal/repository"
	"github.com/madhurima877/chat-system/internal/service"
	"github.com/madhurima877/chat-system/proto/chat"
	"google.golang.org/grpc"
)

func main() {
	mongoClient, err := database.ConnectMongo()
	if err != nil {
		log.Fatalln(err)
	}
	lis, err := net.Listen("tcp", ":50059")
	if err != nil {
		log.Fatalln(err)
	}
	repo := repository.NewChatRepository(mongoClient)
	chatservice := service.NewChatService(repo)

	grpcServer := grpc.NewServer()
	grpcHandler := chatgrpc.NewChatServer(chatservice)
	chat.RegisterChatServiceServer(grpcServer, grpcHandler)

	go func() {

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	grpcServer.GracefulStop()

}
