package control

import (
	"context"
	"net/http"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func CheckoutHandler(c *gin.Context) {
	// 解析请求体
	var req rep.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 构建消息内容
	message := map[string]interface{}{
		"id":           req.ID,
		"user_id":      req.UserID,
		"total":        req.Total,
		"status":       req.Status,
		"created_time": req.CreatedTime,
		"update_time":  req.UpdateTime,
		"order_items":  req.OrderItems,
	}

	// 将消息发送到 Redis 队列
	ctx := context.Background()
	err := redisClient.LPush(ctx, "seckill_queue", message).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push to Redis"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Seckill order submitted"})
}
func CheckoutResultHandler(c *gin.Context) {
	var result struct {
		OrderID int    `json:"order_id"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 根据结果返回响应
	if result.Status == "failed" {
		c.JSON(http.StatusConflict, gin.H{
			"order_id": result.OrderID,
			"status":   result.Status,
			"message":  result.Message,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"order_id": result.OrderID,
			"status":   result.Status,
			"message":  result.Message,
		})
	}
}
