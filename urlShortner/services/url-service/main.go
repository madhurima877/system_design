package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"system_design/urlShortner/internal/cache"
	"system_design/urlShortner/internal/db"
	"system_design/urlShortner/internal/repository"
	"system_design/urlShortner/proto/url"
	"system_design/urlShortner/services/handler"

	"google.golang.org/grpc"
)

func main() {
	redisCache := cache.NewRedisCache()
	if err := redisCache.Ping(); err != nil {
		panic(err)
	}

	log.Println("Redis connected")

	databse, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	lis, err := net.Listen("tcp", ":50061")
	if err != nil {
		panic(err)
	}
	urlRepo := repository.NewURLRepository(databse)
	urlhandler := handler.NewURLHandler(urlRepo, redisCache)

	grpcServer := grpc.NewServer()

	url.RegisterURLServiceServer(grpcServer, urlhandler)

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
