package control

import (
	"fmt"
	"net/http"
	"time"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
)

func CreateOrderHandler(c *gin.Context) {
	var req struct {
		UserID      int     `json:"user_id"`
		Total       float64 `json:"total"`
		ProductID   int     `json:"product_id"`
		Quantity    int     `json:"quantity"`
		Discount    float64 `json:"discount"`
		Payable     float64 `json:"payable"`
		CouponCode  string  `json:"coupon_code"`
		OrderNumber string  `json:"order_number"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 创建订单
	order := rep.Order{
		UserID:      req.UserID,
		Total:       req.Total,
		Status:      "pending",
		CreatedTime: time.Now(),
		UpdateTime:  time.Now(),
		ProductID:   req.ProductID,
		Quantity:    req.Quantity,
		Discount:    req.Discount,
		Payable:     req.Payable,
		CouponCode:  req.CouponCode,
		OrderNumber: req.OrderNumber,
	}
	if err := rep.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// 删除 Redis 缓存中该用户的订单信息
	cacheKey := fmt.Sprintf("user_orders_%d", req.UserID)
	if err := redisClient.Del(ctx, cacheKey).Err(); err != nil {
		fmt.Println("删除 Redis 缓存失败:", err) // 打印删除缓存失败信息
	} else {
		fmt.Println("成功删除 Redis 缓存中的用户订单信息") // 打印删除缓存成功信息
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully"})
}

func GetProductStockHandler(c *gin.Context) {

	// 调用服务层获取商品库存
	products, err := rep.GetAll[rep.Product](rep.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	inventory := make([]gin.H, 0)
	for _, product := range products {
		inventory = append(inventory, gin.H{
			"id":    product.ID,
			"name":  product.Name,
			"stock": product.Stock,
		})
	}

	// 返回库存信息
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   inventory,
	})
}

// 更新库存控制器
func UpdateStockHandler(c *gin.Context) {
	type UpdateStockRequest struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	var request UpdateStockRequest

	// 解析请求体
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 分布式锁的键
	lockKey := fmt.Sprintf("lock:product:%d", request.ProductID)
	lockTimeout := 10 * time.Second // 锁的超时时间

	// 尝试获取分布式锁
	lockAcquired, err := acquireLock(lockKey, lockTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire lock"})
		return
	}
	if !lockAcquired {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to acquire lock, please retry"})
		return
	}

	// 确保在函数结束时释放锁
	defer func() {
		if err := releaseLock(lockKey); err != nil {
			fmt.Println("Failed to release lock:", err)
		}
	}()

	// 查询商品库存
	var product rep.Product
	if err := rep.DB.First(&product, request.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// 计算新的库存
	newStock := product.Stock + request.Quantity
	if newStock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	// 更新库存
	product.Stock = newStock
	if err := rep.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Stock updated successfully",
		"product": product,
	})
}

func RepAPIInit() {
	r := gin.Default()

	r.POST("/orders/create", CreateOrderHandler)
	r.GET("/products/stock", GetProductStockHandler)
	r.POST("/products/stock/update", UpdateStockHandler)

	r.Run(":8081")

	fmt.Println("Rep API Start at 8081 Successfully")
}
