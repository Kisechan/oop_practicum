package control

import (
	"net/http"
	"strconv"
	"time"
	"web_sql/rep"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取单个运单
func GetShipping(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shipping ID"})
		return
	}

	var shipping rep.Shipping
	if err := rep.DB.First(&shipping, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shipping not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, shipping)
}

// 查询所有运单（分页）
func ListShipping(c *gin.Context) {
	var shippings []rep.Shipping
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	if err := rep.DB.Offset(offset).Limit(pageSize).Find(&shippings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shippings"})
		return
	}

	c.JSON(http.StatusOK, shippings)
}

// 创建新运单
func CreateShipping(c *gin.Context) {
	var shipping rep.Shipping
	if err := c.ShouldBindJSON(&shipping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 设置默认创建时间
	now := time.Now()
	shipping.CreateTime = &now

	if err := rep.DB.Create(&shipping).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shipping)
}

// 更新运单信息
func UpdateShipping(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shipping ID"})
		return
	}

	var shipping rep.Shipping
	if err := rep.DB.First(&shipping, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shipping not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var updateData rep.Shipping
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := rep.DB.Model(&shipping).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping updated successfully", "shipping": shipping})
}

// 删除单个运单
func DeleteShipping(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shipping ID"})
		return
	}

	if err := rep.DB.Delete(&rep.Shipping{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping deleted successfully"})
}

// 根据追踪号查询运单状态
func TrackShipping(c *gin.Context) {
	trackingNumber := c.Param("trackingNumber")
	var shipping rep.Shipping

	if err := rep.DB.Where("tracking_number = ?", trackingNumber).First(&shipping).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shipping not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": shipping.Status, "estimated_delivery": shipping.EstimatedDeliveredTime})
}

// 查询某订单项的所有运单
func ListShippingByOrderItem(c *gin.Context) {
	orderItemID, err := strconv.Atoi(c.Param("orderItemID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order item ID"})
		return
	}

	var shippings []rep.Shipping
	if err := rep.DB.Where("order_item_id = ?", orderItemID).Find(&shippings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shippings"})
		return
	}

	c.JSON(http.StatusOK, shippings)
}
