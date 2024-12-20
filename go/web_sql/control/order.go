package control

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

type OrderRequest struct {
	OrderID   string  `json:"order_id"`
	UserID    int     `json:"user_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
}

// func generateOrderID() string {
// 	// 使用时间戳和随机数生成订单号
// 	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
// 	random := rand.Intn(10000) // 生成一个 0-9999 的随机数
// 	return fmt.Sprintf("ORD%d%04d", timestamp, random)
// }

func CheckoutHandler(c *gin.Context) {
	// 解析请求体
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	fmt.Println(req)
	// ctx := context.Background()
	// _, err := redisClient.Ping(ctx).Result()
	// if err != nil {
	// 	log.Fatalf("Failed to connect to Redis: %v", err)
	// }
	// 将秒杀请求发送到 Redis 队列
	ctx := context.Background()
	message, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}
	err = redisClient.LPush(ctx, "seckill_queue", message).Err()
	if err != nil {
		fmt.Printf("Failed to send to Redis: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push to Redis"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Seckill order submitted"})
}

// 秒杀结果的缓存
var seckillResults = struct {
	sync.RWMutex
	results map[string]string
}{results: make(map[string]string)}

// 处理秒杀结果的控制器
func CheckoutResultHandler(c *gin.Context) {
	var result struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 将秒杀结果存入缓存
	seckillResults.Lock()
	seckillResults.results[result.OrderID] = result.Status
	seckillResults.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Seckill result saved"})
}

// 查询秒杀结果的控制器
func GetCheckoutResultHandler(c *gin.Context) {
	orderID := c.Query("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing order_id"})
		return
	}

	// 从缓存中获取秒杀结果
	seckillResults.RLock()
	status, exists := seckillResults.results[orderID]
	seckillResults.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seckill result not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order_id": orderID, "status": status})
}

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
