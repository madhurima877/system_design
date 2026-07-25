package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"notification-system/models"
	pb "notification-system/proto/notification"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NotificationHandler(client pb.NotificationServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Message
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp, err := client.SendNotification(context.Background(), &pb.SendNotificationRequest{
			UserId:  req.UserId,
			Message: req.Message,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)
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
