package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Job struct {
	Id   int
	Data string
}

func main() {
	fmt.Println("Kafka Manual Commit Demo")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:39092"},
		Topic:   "manual-commit-demo",
	})

	defer writer.Close()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:39092"},
		Topic:   "manual-commit-demo",
		GroupID: "manual-commit-group",
	})

	defer reader.Close()

	jobs := make(chan Job)
	var writerWg sync.WaitGroup

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 3; i++ {
		writerWg.Add(1)
		go func(writerWg *sync.WaitGroup) {
			defer writerWg.Done()
			for i := range jobs {
				data, err := json.Marshal(i)
				if err != nil {
					return
				}
				if err := writer.WriteMessages(context, kafka.Message{Value: data}); err != nil {
					log.Println(err)
				}
			}

		}(&writerWg)
	}

	for i := 0; i < 3; i++ {
		writerWg.Add(1)
		go func(writerWg *sync.WaitGroup) {
			defer writerWg.Done()
			for {
				msg, err := reader.FetchMessage(context)
				if err != nil {
					if context.Err() != nil {
						return
					}
					continue
				}

				var e Job
				if err := json.Unmarshal(msg.Value, &e); err != nil {
					log.Println(err)
					continue
				}

				log.Println(e)

				err = reader.CommitMessages(context, msg)
				if err != nil {
					log.Println(err)
				}
			}

		}(&writerWg)
	}
	go func() {
		for i := 0; i <= 100; i++ {
			jobs <- Job{Id: i, Data: "Added Data"}
		}
		close(jobs)
	}()
	writerWg.Wait()
}
