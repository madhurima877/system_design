package main

import (
	"context"
	"log"
	"net/http"
	pb "notification-system/proto/notification"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NotificationHandler(client pb.NotificationServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client.SendNotification(context.Background(), &pb.SendNotificationRequest{
			UserId:  "123",
			Message: "Hi",
		})
	}
}
func main() {
	conn, err := grpc.NewClient("localhost:50057", grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatalln(err)
	}
	client := pb.NewNotificationServiceClient(conn)
	http.HandleFunc("/send/notification", NotificationHandler(client))

	http.ListenAndServe(":8080", nil)
}
