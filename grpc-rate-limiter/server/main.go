package main

import (
	"grpc-rate-limiter/proto/hello"
	"grpc-rate-limiter/server/handler"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50058")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected")
	grpcServer := grpc.NewServer()
	handler := handler.NewServer()

	hello.RegisterGreeterServer(grpcServer, handler)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
