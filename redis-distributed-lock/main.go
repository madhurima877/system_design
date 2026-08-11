package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:36380",
	})
	return &RedisClient{client: client}
}
func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

type DistributedLock struct {
	redis   *RedisClient
	ownerID string
}

func NewDistributedLock(redis *RedisClient, ownerid string) *DistributedLock {
	return &DistributedLock{redis: redis, ownerID: ownerid}
}
func (d *DistributedLock) Acquire(orderID string) bool {
	key := "lock:" + orderID

	acquired := d.redis.client.SetNX(
		context.Background(),
		key,
		d.ownerID,
		10*time.Second,
	).Val()
	if acquired == false {
		return false
	}
	return true
}

var releaseScript = redis.NewScript(`
local value = redis.call("GET", KEYS[1])

if value == ARGV[1] then
    return redis.call("DEL", KEYS[1])
end

return 0
`)

func (d *DistributedLock) Release(orderId string) {
	key := "lock:" + orderId

	releaseScript.Run(
		context.Background(),
		d.redis.client,
		[]string{key},
		d.ownerID,
	)
}
func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recovered from panic")
		}
	}()
	ctx := context.Background()
	redisCache := NewRedisClient()
	if err := redisCache.Ping(ctx); err != nil {
		panic(err)

	}
	// lock := NewDistributedLock(redisCache, "abc-124")
	// acquired1 := lock.Acquire("order-123")
	// log.Println(acquired1)

	// lock.Release("order-123")

	// acquired2 := lock.Acquire("order-123")
	// log.Println(acquired2)

	lock1 := NewDistributedLock(redisCache, "abc-124")
	lock2 := NewDistributedLock(redisCache, "xyz-999")
	log.Println(lock1.Acquire("order-123"))
	log.Println(lock2.Acquire("order-123"))

	lock2.Release("order-123")

	log.Println(lock1.Acquire("order-123"))

}
