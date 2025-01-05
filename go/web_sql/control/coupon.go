package control

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
)

func GetUserCouponsHandler(c *gin.Context) {
	userID := c.Param("id")

	// 尝试从 Redis 缓存中获取优惠券信息
	cacheKey := fmt.Sprintf("user_coupons_%s", userID)
	cachedCoupons, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// 如果缓存中存在优惠券信息，直接返回
		var coupons []rep.Coupon
		if err := json.Unmarshal([]byte(cachedCoupons), &coupons); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"coupons": coupons,
			})
			return
		}
	}

	// 查询用户的优惠券
	var coupons []rep.Coupon
	if err := rep.DB.Where("user_id = ? AND status = ?", userID, "available").Find(&coupons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch coupons"})
		return
	}

	// 将优惠券信息存入 Redis 缓存
	couponsJSON, err := json.Marshal(coupons)
	if err == nil {
		redisClient.Set(ctx, cacheKey, couponsJSON, time.Hour) // 缓存 1 小时
	}

	// 返回优惠券信息
	c.JSON(http.StatusOK, gin.H{
		"coupons": coupons,
	})
}

func UseCouponHandler(c *gin.Context) {
	// 解析请求体
	var request struct {
		UserID     int    `json:"user_id"`
		CouponCode string `json:"coupon_code"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 根据优惠券码和用户ID查询优惠券
	var coupon rep.Coupon
	if err := rep.DB.Where("code = ? AND user_id = ?", request.CouponCode, request.UserID).First(&coupon).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	// 检查优惠券是否已过期
	currentTime := time.Now()
	if currentTime.After(coupon.ExpirationDate) {
		// 更新优惠券状态为 "expired"
		coupon.Status = "expired"
		if err := rep.DB.Save(&coupon).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coupon status"})
			return
		}

		// 返回失败响应
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Coupon has expired",
			"coupon":  coupon,
		})
		return
	}

	// 检查优惠券状态是否为 "available"
	if coupon.Status != "available" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Coupon is not available",
			"coupon":  coupon,
		})
		return
	}

	// 更新优惠券状态为 "used"
	coupon.Status = "used"
	if err := rep.DB.Save(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to use coupon"})
		return
	}

	// 删除与该用户相关的优惠券缓存
	cacheKey := fmt.Sprintf("user_coupons_%d", coupon.UserID)
	redisClient.Del(ctx, cacheKey)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Coupon used successfully",
		"coupon":  coupon,
	})
}
