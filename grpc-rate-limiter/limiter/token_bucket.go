package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate int
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}
func (t *TokenBucket) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	elapsed := time.Since(t.lastRefill)
	seconds := int(elapsed.Seconds())
	tokensToADD := seconds * t.refillRate

	if seconds > 0 {
		t.tokens = min(t.capacity, t.tokens+tokensToADD)
		t.lastRefill = t.lastRefill.Add(time.Duration(seconds) * time.Second)
	}
	if t.tokens > 0 {
		t.tokens -= 1
		return true
	} else {
		return false
	}

}
