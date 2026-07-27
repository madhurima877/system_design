package main

import (
	"context"
	"encoding/json"
	"grpc-rate-limiter/limiter"
	"grpc-rate-limiter/proto/hello"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50058", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	grpcClient := hello.NewGreeterClient(conn)
	bucket := limiter.NewTokenBucket(10, 1)
	http.HandleFunc("/send/request", ClientHandler(grpcClient, bucket))

	http.ListenAndServe(":8080", nil)
}

func ClientHandler(client hello.GreeterClient, bucket *limiter.TokenBucket) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !bucket.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		resp, err := client.SayHello(context.Background(), &hello.GreeterRequest{
			Name: "Madhurima",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
}
