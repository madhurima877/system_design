package workers

import (
	"fmt"
	"system_design/real-time-notification-server/websocket"
	"time"
)

type Job struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}
type WorkerPool struct {
	Jobs    chan Job
	Workers int
	Hub     *websocket.Hub
}

func NewWorkerPool(hub *websocket.Hub) *WorkerPool {
	return &WorkerPool{
		Jobs:    make(chan Job, 10),
		Workers: 3,
		Hub:     hub,
	}
}

func (wp *WorkerPool) Start() {

	for i := 0; i < wp.Workers; i++ {

		go wp.Worker(i)
	}

}

func (wp *WorkerPool) Worker(id int) {

	for job := range wp.Jobs {
		fmt.Printf("Worker %d processing Job %d: %s\n", id, job.ID, job.Message)
		time.Sleep(5 * time.Second)
		wp.Hub.Broadcast([]byte(job.Message))

		fmt.Printf("Worker %d finished\n", id)
	}

}
