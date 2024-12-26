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

func CreateReviewHandler(c *gin.Context) {
	var review rep.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := rep.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	// 删除与该商品相关的评论缓存
	cacheKey := fmt.Sprintf("product_reviews_%d", review.ProductID)
	redisClient.Del(ctx, cacheKey)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"review":  review,
	})
}

func GetProductReviewsHandler(c *gin.Context) {
	productID := c.Param("id")

	// 尝试从 Redis 缓存中获取评论信息
	cacheKey := fmt.Sprintf("product_reviews_%s", productID)
	cachedReviews, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// 如果缓存中存在评论信息，直接返回
		var reviews []rep.Review
		if err := json.Unmarshal([]byte(cachedReviews), &reviews); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"reviews": reviews,
			})
			return
		} else {
			fmt.Println("Redis 缓存解析失败:", err)
		}
	} else {
		fmt.Println("Redis 缓存未命中:", err)
	}

	// 查询商品的评论
	var reviews []rep.Review
	productIDint, err := strconv.Atoi(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	if err := rep.DB.Preload("User").Preload("Product").Where("product_id = ?", productIDint).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	// 将评论信息存入 Redis 缓存
	reviewsJSON, err := json.Marshal(reviews)
	if err == nil {
		redisClient.Set(ctx, cacheKey, reviewsJSON, time.Hour) // 缓存 1 小时
	} else {
		fmt.Println("Redis 缓存写入失败:", err)
	}

	// 返回评论信息
	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
	})
}
