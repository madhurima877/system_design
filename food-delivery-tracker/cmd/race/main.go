package main

import (
	"fmt"
	"sync"
)

var count int
var mu sync.Mutex

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			count++
		}(&wg)
	}
	wg.Wait()

	fmt.Println(count)
}
