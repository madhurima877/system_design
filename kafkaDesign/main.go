package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Message struct {
	ID     int
	UserID string
	Data   string
}
type ProcessFunc func(Message) (string, error)
type RateLimiter struct {
	limit        int
	window       time.Duration
	mu           sync.Mutex
	requestCount map[string]int
	start        map[string]time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:        limit,
		window:       window,
		requestCount: make(map[string]int),
		start:        make(map[string]time.Time),
	}
}
func (rl *RateLimiter) Allow(userID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.start[userID].IsZero() || time.Since(rl.start[userID]) >= rl.window {
		rl.start[userID] = time.Now()
		rl.requestCount[userID] = 0

	}
	if rl.requestCount[userID] >= rl.limit {
		return false
	}
	rl.requestCount[userID]++
	return true
}

type Store struct {
	mu   sync.RWMutex
	data map[int]string
}

func NewStore() *Store {
	return &Store{data: make(map[int]string)}
}

func (s *Store) Save(id int, result string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = result

}
func (s *Store) Count() int {
	return len(s.data)

}

type Processor struct {
	workerCount int
	rateLimiter *RateLimiter
	store       *Store
	processFn   ProcessFunc
	jobs        chan Message
	wg          sync.WaitGroup
	mu          sync.RWMutex
	processed   int
	failed      int
}

func NewProcessor(workerCount int, rl *RateLimiter, store *Store, fn ProcessFunc) *Processor {
	return &Processor{
		workerCount: workerCount,
		rateLimiter: rl,
		store:       store,
		processFn:   fn,
		jobs:        make(chan Message, 100),
	}
}

func (p *Processor) Start() {
	for i := 1; i <= p.workerCount; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for msg := range p.jobs {
				if !p.rateLimiter.Allow(msg.UserID) {
					p.mu.Lock()
					p.failed++
					p.mu.Unlock()
					continue
				}
				success := false

				for attempt := 0; attempt < 3; attempt++ {
					log.Printf("Processing ID %d, attempt %d", msg.ID, attempt+1)
					result, err := p.processFn(msg)
					if err == nil {
						p.store.Save(msg.ID, result)
						p.mu.Lock()
						p.processed++
						p.mu.Unlock()
						success = true
						break
					}
					time.Sleep(time.Duration(1<<attempt) * 100 * time.Millisecond)
				}
				if !success {
					p.mu.Lock()
					p.failed++
					p.mu.Unlock()
				}

			}

		}()
	}
}
func (p *Processor) Submit(msg Message) {
	p.jobs <- msg
}
func (p *Processor) Shutdown() {
	close(p.jobs)
	p.wg.Wait()

}

func (p *Processor) Stats() (int, int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.processed, p.failed

}

func defaultProcess(msg Message) (string, error) {
	time.Sleep(10 * time.Millisecond)
	if msg.ID%7 == 0 {
		return "", errors.New("error")
	}
	return "processed" + msg.Data, nil

}

func main() {
	rl := NewRateLimiter(2, time.Second)
	store := NewStore()
	p := NewProcessor(5, rl, store, defaultProcess)
	p.Start()
	for i := 1; i <= 50; i++ {
		p.Submit(Message{
			ID:     i,
			UserID: fmt.Sprintf("user%d", i%3),
			Data:   fmt.Sprintf("data-%d", i),
		})
	}
	p.Shutdown()
	processed, failed := p.Stats()
	log.Println(processed)
	log.Println(failed)

}
