package main

import (
	"log"
	"net"
	"notification-system/notification-service/handler"
	"notification-system/proto/notification"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50057")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	notihandler := handler.NewNotificationHandler()

	notification.RegisterNotificationServiceServer(grpcServer, notihandler)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
