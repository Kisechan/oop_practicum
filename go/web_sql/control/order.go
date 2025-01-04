package control

import (
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

type CheckoutRequest struct {
	UserID      int     `json:"user_id"`
	ProductID   int     `json:"product_id"`
	Quantity    int     `json:"quantity"`
	CouponCode  string  `json:"coupon_code"`
	Discount    float64 `json:"discount"`
	Payable     float64 `json:"payable"`
	OrderNumber string  `json:"order_number"`
	Total       float64 `json:"total"`
}

func CheckoutHandler(c *gin.Context) {
	var request CheckoutRequest

	// 解析客户端发送的JSON请求
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 构造结算请求
	checkoutRequest := map[string]interface{}{
		"order_number": request.OrderNumber,
		"user_id":      request.UserID,
		"product_id":   request.ProductID,
		"quantity":     request.Quantity,
		"coupon_code":  request.CouponCode,
		"discount":     float64(request.Discount), // 确保是浮点数
		"payable":      float64(request.Payable),  // 确保是浮点数
		"total":        float64(request.Total),    // 确保是浮点数
	}

	// 将结算请求推送到Redis消息队列
	err := PushToRedisQueue("checkout_queue", checkoutRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push to Redis queue"})
		return
	}

	// 返回成功响应
	fmt.Println("成功将JSON: ", checkoutRequest, " 推送到Redis队列")
	c.JSON(http.StatusOK, gin.H{
		"order_number": request.OrderNumber,
		"message":      "Checkout request received",
	})
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
	// 获取订单号
	orderNumber := c.Param("order_number")

	// 从Redis中获取订单结果
	orderResult, err := GetFromRedis("order_result:" + orderNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// 返回订单结果
	c.JSON(http.StatusOK, gin.H{
		"order_number": orderNumber,
		"result":       orderResult,
	})
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
