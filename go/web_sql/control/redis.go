package control

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	err = redisClient.LPush(ctx, "checkout_queue", "").Err()
	if err != nil {
		log.Fatalf("Failed to create Redis queue: %v", err)
	}
	fmt.Println("Connected to Redis and Create Redis Queue Successfully!")
}

// 将结算请求推送到Redis消息队列
func PushToRedisQueue(queueName string, data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return redisClient.LPush(context.Background(), queueName, jsonData).Err()
}

// 从Redis中获取订单结果
func GetFromRedis(key string) (map[string]interface{}, error) {
	val, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// 获取分布式锁
func acquireLock(lockKey string, timeout time.Duration) (bool, error) {
	// 使用 SETNX 命令尝试获取锁
	result, err := redisClient.SetNX(ctx, lockKey, "locked", timeout).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

// 释放分布式锁
func releaseLock(lockKey string) error {
	// 使用 DEL 命令释放锁
	_, err := redisClient.Del(ctx, lockKey).Result()
	return err
}
