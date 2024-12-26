package control

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
)

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

func GetOrdersHandler(c *gin.Context) {
	userID := c.Param("id")
	fmt.Println("获取到的 userID:", userID) // 打印 userID

	// 尝试从 Redis 缓存中获取用户订单信息
	cacheKey := fmt.Sprintf("user_orders_%s", userID)
	fmt.Println("生成的 Redis 缓存键:", cacheKey) // 打印缓存键

	cachedOrders, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		fmt.Println("从 Redis 缓存中获取到用户订单信息") // 打印缓存命中信息

		// 如果缓存中存在用户订单信息，直接返回
		var orders []rep.Order
		if err := json.Unmarshal([]byte(cachedOrders), &orders); err == nil {
			fmt.Println("成功解析 Redis 缓存中的用户订单信息") // 打印解析成功信息
			c.JSON(http.StatusOK, gin.H{
				"orders": orders,
			})
			return
		} else {
			fmt.Println("解析 Redis 缓存中的用户订单信息失败:", err) // 打印解析失败信息
		}
	} else {
		fmt.Println("未从 Redis 缓存中获取到用户订单信息:", err) // 打印缓存未命中信息
	}

	// 查询数据库中的用户订单信息
	var orders []rep.Order
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := rep.DB.Preload("User").Preload("Product").Where("user_id = ?", userIDint).Find(&orders).Error; err != nil {
		fmt.Println("查询数据库中的用户订单信息失败:", err) // 打印查询失败信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	// 将用户订单信息存入 Redis 缓存
	ordersJSON, err := json.Marshal(orders)
	if err == nil {
		fmt.Println("成功将用户订单信息序列化为 JSON") // 打印序列化成功信息
		if err := redisClient.Set(ctx, cacheKey, ordersJSON, time.Hour).Err(); err == nil {
			fmt.Println("成功将用户订单信息存入 Redis 缓存") // 打印缓存成功信息
		} else {
			fmt.Println("将用户订单信息存入 Redis 缓存失败:", err) // 打印缓存失败信息
		}
	} else {
		fmt.Println("将用户订单信息序列化为 JSON 失败:", err) // 打印序列化失败信息
	}

	// 返回用户订单信息
	fmt.Println("返回用户订单信息给客户端") // 打印返回信息
	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}
