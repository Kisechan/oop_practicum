package control

import (
	"net/http"
	"strconv"
	"time"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取单个优惠券
func GetCoupon(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coupon ID"})
		return
	}

	var coupon rep.Coupon
	if err := rep.DB.Preload("User").Preload("Product").First(&coupon, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, coupon)
}

// 查询所有优惠券（分页）
func ListCoupons(c *gin.Context) {
	var coupons []rep.Coupon
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	if err := rep.DB.Offset(offset).Limit(pageSize).Preload("User").Preload("Product").Find(&coupons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch coupons"})
		return
	}

	c.JSON(http.StatusOK, coupons)
}

// 创建新优惠券
func CreateCoupon(c *gin.Context) {
	var coupon rep.Coupon
	if err := c.ShouldBindJSON(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if coupon.ExpirationDate != nil && coupon.ExpirationDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expiration date cannot be in the past"})
		return
	}

	if err := rep.DB.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, coupon)
}

// 更新优惠券
func UpdateCoupon(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coupon ID"})
		return
	}

	var coupon rep.Coupon
	if err := rep.DB.First(&coupon, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var updateData rep.Coupon
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 更新数据
	coupon.Code = updateData.Code
	coupon.Type = updateData.Type
	coupon.Discount = updateData.Discount
	coupon.Minimum = updateData.Minimum
	coupon.ExpirationDate = updateData.ExpirationDate
	coupon.Status = updateData.Status

	if err := rep.DB.Save(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coupon updated successfully", "coupon": coupon})
}

// 删除优惠券
func DeleteCoupon(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coupon ID"})
		return
	}

	if err := rep.DB.Delete(&rep.Coupon{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coupon deleted successfully"})
}

// 查询某用户的优惠券
func GetCouponsByUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var coupons []rep.Coupon
	if err := rep.DB.Where("user_id = ?", userID).Preload("Product").Find(&coupons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user's coupons"})
		return
	}

	c.JSON(http.StatusOK, coupons)
}

// 根据优惠码检查优惠券
func CheckCouponByCode(c *gin.Context) {
	code := c.Param("code")

	var coupon rep.Coupon
	if err := rep.DB.Where("code = ?", code).Preload("User").Preload("Product").First(&coupon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 检查是否过期或已使用
	if coupon.Status == "expired" || coupon.ExpirationDate != nil && coupon.ExpirationDate.Before(time.Now()) {
		coupon.Status = "expired"
		rep.DB.Save(&coupon)
		c.JSON(http.StatusOK, gin.H{"message": "Coupon is expired", "coupon": coupon})
		return
	}

	if coupon.Status == "used" {
		c.JSON(http.StatusOK, gin.H{"message": "Coupon has already been used", "coupon": coupon})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coupon is valid", "coupon": coupon})
}
