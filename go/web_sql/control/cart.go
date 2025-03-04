package control

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
)

func AddCartHandler(c *gin.Context) {
	var cartItem rep.Cart

	// 绑定 JSON 数据到 cartItem 结构体
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否存在
	var user rep.User
	if err := rep.DB.First(&user, cartItem.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 检查产品是否存在
	var product rep.Product
	if err := rep.DB.First(&product, cartItem.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// 检查库存是否足够
	// if product.Stock < cartItem.Quantity {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
	// 	return
	// }

	// 创建购物车项
	cartItem.AddTime = time.Now()
	if err := rep.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	// 删除与该用户相关的购物车缓存
	cacheKey := fmt.Sprintf("user_cart_%d", cartItem.UserID)
	redisClient.Del(ctx, cacheKey)

	// 返回成功响应
	c.JSON(http.StatusCreated, gin.H{
		"message": "Item added to cart successfully",
		"cart":    cartItem,
	})
}

func RemoveCartHandler(c *gin.Context) {
	// 获取路径参数 id
	cartItemID := c.Param("id")

	// 查询购物车项
	var cartItem rep.Cart
	if err := rep.DB.First(&cartItem, cartItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	// 删除购物车项
	if err := rep.DB.Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	// 删除与该用户相关的购物车缓存
	cacheKey := fmt.Sprintf("user_cart_%d", cartItem.UserID)
	redisClient.Del(ctx, cacheKey)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from cart successfully",
	})
}

func GetCartHandler(c *gin.Context) {
	userID := c.Query("user_id")

	// 尝试从 Redis 缓存中获取购物车信息
	cacheKey := fmt.Sprintf("user_cart_%s", userID)

	cachedCartItems, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		fmt.Println("从 Redis 缓存中获取到购物车信息") // 打印缓存命中信息

		// 如果缓存中存在购物车信息，直接返回
		var cartItems []rep.Cart
		if err := json.Unmarshal([]byte(cachedCartItems), &cartItems); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"cart_items": cartItems,
			})
			return
		} else {
			fmt.Println("解析 Redis 缓存中的购物车信息失败:", err) // 打印解析失败信息
		}
	} else {
		fmt.Println("未从 Redis 缓存中获取到购物车信息:", err) // 打印缓存未命中信息
	}

	// 查询用户的购物车项
	var cartItems []rep.Cart
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := rep.DB.Preload("Product").Where("user_id = ?", userIDint).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	// 将购物车信息存入 Redis 缓存
	cartItemsJSON, err := json.Marshal(cartItems)
	if err == nil {
		fmt.Println("成功将购物车信息序列化为 JSON") // 打印序列化成功信息
		if err := redisClient.Set(ctx, cacheKey, cartItemsJSON, time.Hour).Err(); err == nil {
			fmt.Println("成功将购物车信息存入 Redis 缓存") // 打印缓存成功信息
		} else {
			fmt.Println("将购物车信息存入 Redis 缓存失败:", err) // 打印缓存失败信息
		}
	} else {
		fmt.Println("将购物车信息序列化为 JSON 失败:", err) // 打印序列化失败信息
	}

	// 返回购物车信息
	c.JSON(http.StatusOK, gin.H{
		"cart_items": cartItems,
	})
}
