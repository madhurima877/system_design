package main

import (
	"log"
	"net/http"
	"system_design/urlShortner/api-gateway/handler"
	"system_design/urlShortner/proto/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50061", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	grpcClient := url.NewURLServiceClient(conn)
	http.HandleFunc("/create/short/url", handler.CreateShortURL(grpcClient))
	http.HandleFunc("/get/original/url", handler.GetOriginalURL(grpcClient))
	http.ListenAndServe(":8080", nil)

}
