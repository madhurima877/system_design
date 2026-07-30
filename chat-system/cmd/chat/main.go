package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	chatgrpc "github.com/madhurima877/chat-system/internal/grpc"
	"github.com/madhurima877/chat-system/proto/chat"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50059")
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	grpcHandler := chatgrpc.NewChatServer()
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
