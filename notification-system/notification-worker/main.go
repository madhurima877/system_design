package main

import (
	"context"
	"log"
	"notification-system/notification-worker/kafka"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	consumer := kafka.NewConsumer()

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go consumer.Worker(ctx, &wg)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	cancel()
	wg.Wait()
	if err := consumer.Close(); err != nil {
		log.Println("failed to close kafka consumer:", err)
	}
}
