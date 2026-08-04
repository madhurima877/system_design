package handlers

import (
	"fmt"
	"sync"
	"time"
)

func ProcessNotification(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Processing notification...")

	for i := 1; i <= 5; i++ {
		fmt.Println("Step", i)
		time.Sleep(time.Second)
	}

	fmt.Println("Notification processed")
}
