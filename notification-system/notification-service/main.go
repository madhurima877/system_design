package main

import (
	"log"
	"net"
	"notification-system/notification-service/handler"
	"notification-system/notification-service/kafka"
	"notification-system/proto/notification"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50057")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	producer := kafka.NewProducer()
	notihandler := handler.NewNotificationHandler(producer)

	notification.RegisterNotificationServiceServer(grpcServer, notihandler)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	grpcServer.GracefulStop()
	if err := producer.Close(); err != nil {
		log.Println("failed to close kafka producer:", err)
	}

}
