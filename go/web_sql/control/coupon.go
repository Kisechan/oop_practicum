package control

import (
	"net/http"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
)

func GetUserCouponsHandler(c *gin.Context) {
	userID := c.Param("id")

	var coupons []rep.Coupon
	if err := rep.DB.Where("user_id = ?", userID).Find(&coupons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coupons": coupons,
	})
}

func UseCouponHandler(c *gin.Context) {
	var coupon rep.Coupon
	if err := c.ShouldBindJSON(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coupon.Status = "used"
	if err := rep.DB.Save(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to use coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Coupon used successfully",
		"coupon":  coupon,
	})
}
