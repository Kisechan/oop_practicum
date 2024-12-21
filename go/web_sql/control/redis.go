package control

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()

	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
)

func RedisInit() {
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 创建队列
	err = redisClient.LPush(ctx, "seckill_queue", "").Err()
	if err != nil {
		log.Fatalf("Failed to create Redis queue: %v", err)
	}
	fmt.Println("Connected to Redis and Create Redis Queue Successfully!")
}
